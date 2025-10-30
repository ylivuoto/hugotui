package main

import (
	"strings"

	"hugotui/commands"
	"hugotui/utils"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func (m model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Bold(true).
			Padding(0, 1, 0, 2).Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("62")),
	)
}

func (m *model) View() string {
	if m.focus == 1 {
		m.viewport.Style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	}

	if m.focus == 2 {
		return m.handleEditArticle()
	}
	if m.focus == 3 {
		return m.handleCreateArticle()
	}

	if !m.ready {
		return "\n  Initializing..."
	}
	return m.showView("default")
}

func (m *model) handleCreateArticle() string {
	completed := m.form.State == huh.StateCompleted
	confirmed := m.form.GetBool("confirm")

	if completed && confirmed {
		heading := m.form.GetString("heading")
		tags := m.form.Get("tags").([]string)

		commands.CreateArticle(heading, tags)
		m.refreshList()
		return m.showView("default")
	}

	if completed && !confirmed {
		m.focus = 0
		return m.showView("default")
	}
	return m.showView("createArticle")
}

func (m *model) handleEditArticle() string {
	if m.form.State == huh.StateCompleted {
		title := m.form.GetString("heading")
		filepath := m.list.SelectedItem().(item).path

		// TODO: handle error
		utils.ModifyFileTitle(filepath, title)
		utils.ModifyFilePath(filepath, title)

		// FIX: changing focus won't work just like that, maybe trigger update?
		m.focus = 0
		m.refreshList()
		return m.showView("default")
	}

	return m.showView("editArticle")
}

func (m *model) refreshList() {
	items := fetchItems()
	list := setupList(items, m.list.Width(), m.list.Height())
	m.list = list
}

func (m *model) showView(view string) string {
	// TODO: refactor styles
	listBoxStyle := lipgloss.NewStyle().
		MarginRight(1).
		Border(lipgloss.RoundedBorder())

	if m.focus == 0 {
		listBoxStyle = listBoxStyle.BorderForeground(lipgloss.Color("62")) // green-ish
	}

	commadLog := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 1).
		Width(m.width - 4).
		Height(8).
		Render(m.cmdLog.View())

	v := strings.TrimSuffix(m.form.View(), "\n\n")
	form := lipgloss.NewStyle().Margin(1, 0).Render(v)
	row := lipgloss.JoinHorizontal(lipgloss.Top, listBoxStyle.Render(m.list.View()), m.viewport.View())
	helpView := m.help.View(m.keys)

	switch view {
	case "createArticle":
		header := m.appBoundaryView("Create article")
		return docStyle.Render(header + "\n" + form)

	case "editArticle":
		header := m.appBoundaryView("Edit article")
		return docStyle.Render(header + "\n" + form)

	default:
		return docStyle.Render(lipgloss.JoinVertical(0, row, commadLog, helpView))
	}
}
