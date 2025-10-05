package main

import (
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
	title, body string
}

type Matter struct {
	Title string   `yaml:"title"`
	Tags  []string `yaml:"tags"`
	Date  string   `yaml:"date"`
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.body }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	viewport viewport.Model
	width    int
	height   int
	focus    int
}

func mainModel() (*model, error) {
	items := []list.Item{
		item{title: "Raspberry Pi’s", body: "I have ’em all over my house"},
		item{title: "Nutella", body: "It's good on toast"},
		item{title: "Bitter melon", body: "It cools you down"},
		item{title: "Nice socks", body: "And by that I mean socks without holes"},
		item{title: "Eight hours of sleep", body: "I had this once"},
		item{title: "Cats", body: "Usually"},
		item{title: "Plantasia, the album", body: "My plants love it too"},
		item{title: "Pour over coffee", body: "It takes forever to make though"},
		item{title: "VR", body: "Virtual reality...what is there to say?"},
		item{title: "Noguchi Lamps", body: "Such pleasing organic forms"},
		item{title: "Linux", body: "Pretty much the best OS"},
		item{title: "Business school", body: "Just kidding"},
		item{title: "Pottery", body: "Wet clay is a great feeling"},
		item{title: "Shampoo", body: "Nothing like clean hair"},
		item{title: "Table tennis", body: "It’s surprisingly exhausting"},
		item{title: "Milk crates", body: "Great for packing in your extra stuff"},
		item{title: "Afternoon tea", body: "Especially the tea sandwich part"},
		item{title: "Stickers", body: "The thicker the vinyl the better"},
		item{title: "20° Weather", body: "Celsius, not Fahrenheit"},
		item{title: "Warm light", body: "Like around 2700 Kelvin"},
		item{title: "The vernal equinox", body: "The autumnal equinox is pretty good too"},
		item{title: "Gaffer’s tape", body: "Basically sticky fabric"},
		item{title: "Terrycloth", body: "In other words, towel fabric"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "My Fave Things"

	const width = 78

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
