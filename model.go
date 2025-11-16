package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"hugotui/commands"
	"hugotui/utils"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(0, 1)

type (
	lineMsg string
	doneMsg struct{ err error }
)

type item struct {
	title, date string
	tags        []string
	content     string
	path        string
}

func (i item) Title() string       { return i.title }
func (i item) FilterValue() string { return i.title }
func (i item) Description() string {
	return fmtDesc(i.date, i.tags)
}

func fmtDesc(dateStr string, tags []string) string {
	return fmt.Sprintf("%s - %s", utils.FormatHugoDate(dateStr), utils.ParseTags(tags))
}

type model struct {
	list     list.Model
	viewport viewport.Model
	renderer *glamour.TermRenderer
	width    int
	height   int
	focus    int
	content  string
	form     *huh.Form
	ready    bool
	cmdLog   viewport.Model
	lines    []string
	done     bool
	sub      chan string
	keys     KeyMap
	help     help.Model
}

func fetchItems() []list.Item {
	posts, _ := commands.ListHugoPosts()
	items := make([]list.Item, len(posts))
	for i, p := range posts {
		items[i] = item{title: p.Title, date: p.Date, tags: p.Tags, content: p.Content, path: p.Path}
	}
	return items
}

func setupList(items []list.Item, width int, heigth int) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), width, heigth)
	l.Title = "All Posts"
	l.SetShowHelp(false)
	return l
}

func mainModel() (*model, error) {
	// Pick all posts via hugo cli
	sub := make(chan string)

	items := fetchItems()
	l := setupList(items, 10, 20)

	const width = 90 // Configure viewport for markdown rendering for Glamour
	vp := viewport.New(width, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder())

	form := newCreateForm([]string{"my", "awesome", "tags"})

	// Help
	help := help.New()
	help.Styles.ShortKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Faint(true).
		PaddingRight(1) // 👈 adds a single space between key & desc

	help.Styles.FullKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Faint(true).
		PaddingRight(2) // slightly more for full view

	return &model{
		list:     l,
		viewport: vp,
		form:     form,
		sub:      sub,
		keys:     Keys,
		help:     help,
	}, nil
}

func (m *model) Init() tea.Cmd {
	m.updateRenderer()
	return tea.Batch(
		m.form.Init(),
		waitForLine(m.sub),
	)
}

func (m *model) updateRenderer() {
	go func() {
		const glamourGutter = 2
		renderer, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(m.viewport.Width-m.viewport.Style.GetHorizontalFrameSize()-glamourGutter),
		)
		m.renderer = renderer
		content := ""
		selectedItem := m.list.SelectedItem()
		if selectedItem != nil {
			content = selectedItem.(item).content
		}

		str, _ := renderer.Render(content)
		m.viewport.SetContent(str)
	}()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.focus < 2 {
			mainViewKeybindings(m, &msg)
		}
		if m.focus == 0 {
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
			m.renderSelected()
		}
		if m.focus == 1 {
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
		if m.focus >= 2 {
			return updateCreate(msg, m)
		}
		// TODO: Handle keyboard and mouse events in the viewport
		// m.cmdLog, cmd = m.cmdLog.Update(msg)
		// cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	// Helps to render properly on window resize
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.width = msg.Width
		m.height = msg.Height
		listContentWidth := min(int(float64(msg.Width)*0.40), 50)
		viewportContentWidth := min(msg.Width-h-listContentWidth, 90)
		m.list.SetSize(listContentWidth, msg.Height-v-11)
		m.viewport.Height = msg.Height - v - 9
		m.viewport.Width = viewportContentWidth - h
		m.viewport.Style.Width(viewportContentWidth - h)
		verticalMarginHeight := m.viewport.Height
		m.help.Width = msg.Width
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.cmdLog = viewport.New(msg.Width, msg.Height-verticalMarginHeight-4)
			m.cmdLog.YPosition = verticalMarginHeight
			m.cmdLog.SetContent(fmt.Sprintf("Test %s", m.content))
			m.ready = true
		} else {
			m.cmdLog.Width = min(msg.Width, listContentWidth+viewportContentWidth-2*h)
			m.cmdLog.Height = msg.Height - verticalMarginHeight - 5
		}

	case lineMsg:
		m.lines = append(m.lines, string(msg))
		m.updateViewport()
		return m, waitForLine(m.sub)

	case doneMsg:
		m.done = true
		m.updateViewport()
	default:
		// Prcess huh form internal messages
		if m.focus >= 2 {
			return updateCreate(msg, m)
		}
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func mainViewKeybindings(m *model, msg *tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Edit):
		m.form = newModifyForm()
		m.form.Init()
		m.focus = 2
	case key.Matches(msg, m.keys.New):
		m.focus = 3
		return updateCreate(msg, m)
	case key.Matches(msg, m.keys.Open):
		utils.OpenFileInEditor(m.list.SelectedItem().(item).path)
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.WordWrap):
		m.updateRenderer()
	case key.Matches(msg, m.keys.Preview):
		// TODO: refactor to own functions
		m.sub <- "Starting Hugo preview server..."
		commands.Preview()
	case key.Matches(msg, m.keys.StopPreview):
		commands.StopPreview()
		m.sub <- "Hugo preview server stopped."
	case key.Matches(msg, m.keys.Push):
		go transferFiles(m.sub)
	case key.Matches(msg, m.keys.Tab):
		if m.focus != 3 {
			m.focus = (m.focus + 1) % 2
		}
	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll
	}
	return m, nil
}

func (m *model) renderSelected() {
	sel := m.list.SelectedItem()
	if sel == nil {
		return
	}
	item := sel.(item)
	out, err := m.renderer.Render(item.content)
	// on render error, show raw body
	if err != nil {
		m.viewport.SetContent(item.content)
		return
	}
	m.viewport.SetContent(out)
	m.viewport.GotoTop()
}

func newCreateForm(tags []string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(

			// TODO: make tags dynamic, pick from existing tags
			huh.NewMultiSelect[string]().
				Key("tags").
				Title("Select tags").
				Options(huh.NewOptions(tags...)...),

			huh.NewInput().
				Key("heading").
				Title("Heading"),

			// TODO: datepicker, use community package
			huh.NewInput().
				Key("date").
				Title("Date"),

			huh.NewConfirm().
				Key("confirm").
				Title("Create article?"),
		),
	)
}

func newModifyForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("heading").
				Title("Title"),
		),
		// TODO: add checkbox for optionally modify file path based on title
	)
}

// Process the form
func updateCreate(msg tea.Msg, m *model) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)

	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}
	return m, cmd
}

func (m *model) updateViewport() {
	content := strings.Join(m.lines, "\n")
	if m.done {
		content += "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Render("✓ Done! ")
	}
	m.cmdLog.SetContent(content)
	m.cmdLog.GotoBottom()
}

func waitForLine(sub chan string) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-sub
		if !ok {
			return doneMsg{}
		}
		return lineMsg(line)
	}
}

func transferFiles(ch chan string) {
	// TODO: move this function
	defer close(ch)
	localDir := "public/"
	remoteDest := utils.HugoRemote
	remoteDir := utils.HugoRemoteDir

	// List all files first
	ch <- "Scanning files..."
	var files []string
	var totalSize int64

	filepath.WalkDir(localDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		files = append(files, path)
		info, _ := d.Info()
		totalSize += info.Size()
		return nil
	})
	ch <- fmt.Sprintf("Found %d files (%.2f MB)\n", len(files), float64(totalSize)/(1024*1024))
	ch <- "Simulating file transfer..."
	// Show each file
	for _, file := range files {
		info, _ := os.Stat(file)
		relPath, _ := filepath.Rel(localDir, file)
		size := info.Size()

		var sizeStr string
		if size < 1024 {
			sizeStr = fmt.Sprintf("%d B", size)
		} else if size < 1024*1024 {
			sizeStr = fmt.Sprintf("%.2f KB", float64(size)/1024)
		} else {
			sizeStr = fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
		}
		ch <- fmt.Sprintf("  %s (%s)", relPath, sizeStr)

	}

	ch <- "\nStarting transfer..."

	// Run scp
	cmd := exec.Command("hugo")
	// cmd := exec.Command("scp", "-r", "-P", port, "public/*", remoteDest)
	out, err := cmd.Output()
	ch <- string(out)
	if err != nil {
		ch <- fmt.Sprintf("\n❌ Error: %v", err)
	} else {
		ch <- "\n Hugo built site!"
	}

	cmd = exec.Command("bash", "-c", fmt.Sprintf("scp -r %s* %s:%s", localDir, remoteDest, remoteDir))
	out, err = cmd.Output()
	ch <- string(out)
	if err != nil {
		ch <- fmt.Sprintf("\n❌ Error: %v", err)
	} else {
		ch <- "\n✅ All files transferred successfully!"
	}
}
