package main

import "github.com/charmbracelet/lipgloss"

func (m model) View() string {
	listBoxStyle := lipgloss.NewStyle().
		MarginRight(1).
		Border(lipgloss.RoundedBorder())
		// Render(m.list.View())

	if m.focus == 0 {
		listBoxStyle = listBoxStyle.BorderForeground(lipgloss.Color("62")) // green-ish
	}

	if m.focus == 1 {
		m.viewport.Style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, listBoxStyle.Render(m.list.View()), m.viewport.View())
	return docStyle.Render(row)
}
