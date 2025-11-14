package main

import (
	"fmt"
	"os"

	"hugotui/commands"
	"hugotui/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	utils.SetupConfig()
	model, err := mainModel()
	if err != nil {
		fmt.Println("Could not initialize Bubble Tea model:", err)
		os.Exit(1)
	}
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Bummer, there's been an error:", err)
		os.Exit(1)
	}
	// NOTE: this might not be called on exit in some cases
	commands.StopPreview()
}
