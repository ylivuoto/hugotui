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

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "My Awesome Posts"

	const width = 78

	// Configure viewport for markdown rendering for Glamour
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// TODO: keybindings help
		// TODO: return to preious view from create
		switch msg.String() {
		case "o":
			utils.OpenFileInEditor(m.list.SelectedItem().(item).path)
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "m":
			m.form = newModifyForm()
			m.form.Init()
			m.focus = 2
		case "tab":
			// Switch focus between list and viewport, but do not switch on create from
			if m.focus != 3 {
				m.focus = (m.focus + 1) % 2
			}
		case "n":
			// TODO: proper keybindings for create new article
			m.focus = 3
			return updateCreate(msg, m)
		default:
			if m.focus == 0 {
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				m.renderSelected()
				return m, cmd
			}
			if m.focus == 1 {
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
			if m.focus >= 2 {
				return updateCreate(msg, m)
			}
		}
	case tea.WindowSizeMsg:
		// Needed for rendering list properly, glamour should work without it
		h, v := docStyle.GetFrameSize()
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-h, msg.Height-v-2)
		m.viewport.Height = msg.Height - v - 2
	default:
		// Prcess huh form internal messages
		if m.focus >= 2 {
			return updateCreate(msg, m)
		}
		return m, nil
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
