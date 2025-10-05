package main

import (
	"fmt"

	"hugotui/commands"
	"hugotui/utils"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, date string
	tags        []string
	content     string
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
}

func mainModel() (*model, error) {
	// Pick all posts via hugo cli
	posts, _ := commands.ListHugoPosts()

	// TODO: fix latest post tags to not to show earlirer posts tags
	// Make bubbletea list items
	items := make([]list.Item, len(posts))
	for i, p := range posts {
		items[i] = item{title: p.Title, date: p.Date, tags: p.Tags, content: p.Content}
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

	str, err := renderer.Render(l.SelectedItem().(item).content)
	if err != nil {
		return nil, err
	}

	vp.SetContent(str)

	return &model{
		list:     l,
		viewport: vp,
		renderer: renderer,
	}, nil
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "tab":
			// Switch focus between list and viewport
			m.focus = (m.focus + 1) % 2
		case "n":
			commands.CreateArticle()
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
				m.renderSelected()
				return m, cmd
			}
		}
	case tea.WindowSizeMsg:
		// Needed for rendering list properly, glamour should work without it
		h, v := docStyle.GetFrameSize()
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-h, msg.Height-v-2)
		m.viewport.Height = msg.Height - v - 2
		m.renderSelected()
	default:
		return m, nil
	}
	return m, nil
}

func (m *model) renderSelected() {
	sel := m.list.SelectedItem()
	item := sel.(item)
	out, err := m.renderer.Render(item.content)
	if err != nil {
		// on render error, show raw body
		m.viewport.SetContent(item.content)
		return
	}
	m.viewport.SetContent(out)
	m.viewport.GotoTop()
}
