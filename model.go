package main

import (
	"fmt"

	"hugotui/commands"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const content = `
# Today’s Menu

## Appetizers

| Name        | Price | Notes                           |
| ---         | ---   | ---                             |
| Tsukemono   | $2    | Just an appetizer               |
| Tomato Soup | $4    | Made with San Marzano tomatoes  |
| Okonomiyaki | $4    | Takes a few minutes to make     |
| Curry       | $3    | We can add squash if you’d like |
`

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, date string
}

// Title as is, join rest of info for desc
func (i item) Title() string       { return i.title }
func (i item) Description() string { return fmt.Sprintf("%s", i.date) }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	viewport viewport.Model
	width    int
	height   int
	focus    int
}

func mainModel() (*model, error) {
	// Pick all posts via hugo cli
	posts, _ := commands.ListHugoPosts()

	// Make bubbletea list items
	items := make([]list.Item, len(posts))
	for i, p := range posts {
		items[i] = item{title: p.Title, date: p.Date}
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

	str, err := renderer.Render(content)
	if err != nil {
		return nil, err
	}

	vp.SetContent(str)

	return &model{
		list:     l,
		viewport: vp,
	}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "tab":
			// Switch focus between list and viewport
			m.focus = (m.focus + 1) % 2
		default:
			if m.focus == 0 {
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			}
			if m.focus == 1 {
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
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
	default:
		return m, nil
	}
	return m, nil
}
