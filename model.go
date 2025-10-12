package main

import (
	"fmt"

	"hugotui/commands"
	"hugotui/utils"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, date string
	tags        []string
	content     string
	path        string
}

// Title as is, join rest of info for desc
func (i item) Title() string { return i.title }

// TODO: fix date format to show day before month
func (i item) Description() string {
	return fmt.Sprintf("%s - %s", utils.FormatHugoDate(i.date), utils.ParseTags(i.tags))
}
func (i item) FilterValue() string { return i.title }

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
}

func mainModel() (*model, error) {
	// Pick all posts via hugo cli
	posts, _ := commands.ListHugoPosts()

	// TODO: fix latest post tags to not to show earlirer posts tags
	// Make bubbletea list items
	items := make([]list.Item, len(posts))
	for i, p := range posts {
		items[i] = item{title: p.Title, date: p.Date, tags: p.Tags, content: p.Content, path: p.Path}
	}

	l := list.New(items, list.NewDefaultDelegate(), 15, 0)
	l.Title = "My Awesome Posts"

	const width = 77 // Configure viewport for markdown rendering for Glamour
	vp := viewport.New(width, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		PaddingRight(2)

	const glamourGutter = 2
	glamourRenderWidth := width - vp.Style.GetHorizontalFrameSize() - glamourGutter

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth),
	)
	if err != nil {
		return nil, err
	}

	content := ""
	selectedItem := l.SelectedItem()
	if selectedItem != nil {
		content = selectedItem.(item).content
	}

	str, err := renderer.Render(content)
	if err != nil {
		fmt.Println("Error rendering markdown:", err)
		return nil, err
	}

	vp.SetContent(str)

	form := newCreateForm([]string{"my", "awesome", "tags"})

	return &model{
		list:     l,
		viewport: vp,
		renderer: renderer,
		form:     form,
	}, nil
}

func (m *model) Init() tea.Cmd {
	return m.form.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	// TODO: implement command logging
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// TODO: keybindings help
		// TODO: return to preious view from create
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
		// TODO: move into separete "focus" scope
		//
		// Handle keyboard and mouse events in the viewport
		// m.cmdLog, cmd = m.cmdLog.Update(msg)
		// cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case tea.WindowSizeMsg:
		// Needed for rendering list properly, glamour should work without it
		h, v := docStyle.GetFrameSize()
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-h, msg.Height-v-12)
		// TODO: fix viewport width on resize, now hardcoded 77 and 60
		// Also height could be relationally sized
		m.viewport.Height = msg.Height - v - 10
		m.viewport.Width = msg.Width - h - 60
		verticalMarginHeight := m.viewport.Height
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
			m.cmdLog.Width = msg.Width
			m.cmdLog.Height = msg.Height - verticalMarginHeight - 5
		}

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
	switch msg.String() {
	case "m":
		m.form = newModifyForm()
		m.form.Init()
		m.focus = 2
	case "n":
		// TODO: proper keybindings for create new article
		m.focus = 3
		return updateCreate(msg, m)
	case "o":
		utils.OpenFileInEditor(m.list.SelectedItem().(item).path)
	case "q", "ctrl+c", "esc":
		return m, tea.Quit
	case "p":
		commands.Preview()
	case "P":
		commands.Publish()
	case "tab":
		// Switch focus between list and viewport, but do not switch on create from
		if m.focus != 3 {
			m.focus = (m.focus + 1) % 2
		}
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
	m.cmdLog.SetContent(item.path)
	m.cmdLog.GotoBottom()
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
			// TODO: confirm button to submit the form
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

// func updateModify(msg tea.Msg, m *model) (tea.Model, tea.Cmd) {
// 	return m, cmd
// }
