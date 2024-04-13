package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type ListKeyMapT struct {
	Next   key.Binding
	Prev   key.Binding
	Open   key.Binding
	New    key.Binding
	Delete key.Binding
    Select key.Binding
	Exit   key.Binding
}

type NoteKeyMapT struct {
	Next   key.Binding
	Prev   key.Binding
	Save   key.Binding
	Cancel key.Binding
}

var ListKeyMap = ListKeyMapT{
	Next: key.NewBinding(
		key.WithKeys("tab", "j", "down"),
		key.WithHelp("↓/j/tab", "move down"),
	),
	Prev: key.NewBinding(
		key.WithKeys("shift+tab", "k", "up"),
		key.WithHelp("↑/k/shift+tab", "move up"),
	),
	Open: key.NewBinding(
		key.WithKeys("enter", "l", "right"),
		key.WithHelp("→/l/enter", "open"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
    Select: key.NewBinding(
        key.WithKeys("v", "shift+v", "ctrl+v", "ctrl+shift+v"),
        key.WithHelp("v", "enter select mode"),
    ),
	Exit: key.NewBinding(
		key.WithKeys("ctrl+c", "q", "escape"),
		key.WithHelp("ctrl+c/q/escape", "exit"),
	),
}

var NoteKeyMap = NoteKeyMapT{
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
	Prev: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "previous field"),
	),
	Save: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open note"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("ctrl+c", "escape"),
		key.WithHelp("", ""),
	),
}

func (k NoteKeyMapT) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Next, k.Prev, k.Save, k.Cancel}}
}

func (k ListKeyMapT) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Prev, k.Exit},
		{k.Open, k.New, k.Delete},
	}
}

func (k NoteKeyMapT) ShortHelp() []key.Binding {
	return []key.Binding{k.Next, k.Prev, k.Save, k.Cancel}
}

func (k ListKeyMapT) ShortHelp() []key.Binding {
	return []key.Binding{k.Open, k.New, k.Delete}
}
