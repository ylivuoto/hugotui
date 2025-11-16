package main

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type KeyMap struct {
	// Navigation
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	// Actions
	Edit        key.Binding
	Delete      key.Binding
	New         key.Binding
	Open        key.Binding
	Tab         key.Binding
	Preview     key.Binding
	Push        key.Binding
	StopPreview key.Binding

	// Common
	Help     key.Binding
	Quit     key.Binding
	WordWrap key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},                                            // first column
		{k.Edit, k.Delete, k.New, k.Open, k.Tab, k.Preview, k.StopPreview, k.Push}, // second column
		{k.Help, k.Quit, k.WordWrap},                                               // third column
	}
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit item"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new item"),
	),
	Delete: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("d", "delete item"),
	),
	Open: key.NewBinding(
		key.WithKeys("o", "enter"),
		key.WithHelp("o", "open item"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "toggle focus"),
	),
	Preview: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "preview item"),
	),
	Push: key.NewBinding(
		key.WithKeys("P"),
		key.WithHelp("P", "push changes"),
	),
	StopPreview: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "stop preview"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	WordWrap: key.NewBinding(key.WithKeys("w"),
		key.WithHelp("w", "toggle word wrap"),
	),
}
