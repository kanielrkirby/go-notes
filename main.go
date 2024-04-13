package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	li "github.com/charmbracelet/lipgloss"
)

var c = map[string]li.Color{
	"primary":          li.Color("#F4442E"),
	"primary-light":    li.Color("#FC9E4F"),
	"primary-lighter":  li.Color("#EDD382"),
	"primary-lightest": li.Color("#F2F3AE"),
	"black":            li.Color("#020122"),
	"white":            li.Color("#FEFEF7"),
}

var (
	listItemTitle            = li.NewStyle().MarginLeft(1).PaddingLeft(1).Bold(true).Foreground(c["white"]).Border(li.HiddenBorder(), false, false, false, true)
	listItemSubtitle         = li.NewStyle().MarginLeft(1).MarginBottom(1).PaddingLeft(1).Foreground(c["white"]).Faint(true).Border(li.HiddenBorder(), false, false, false, true)
	listItemSelectedTitle    = li.NewStyle().MarginLeft(1).PaddingLeft(1).Bold(true).Foreground(c["primary-light"]).Border(li.ThickBorder(), false, false, false, true).BorderForeground(c["primary-light"])
	listItemSelectedSubtitle = li.NewStyle().MarginLeft(1).MarginBottom(1).PaddingLeft(1).Foreground(c["primary-light"]).Faint(true).Border(li.ThickBorder(), false, false, false, true).BorderForeground(c["primary-light"])
)

type Choice struct {
	Title    string
	Subtitle string
}

type Model struct {
	viewport    viewport.Model
	textinput   textinput.Model
	textarea    textarea.Model
	choices     []Choice
	choiceIndex int
	noteIndex   int
	help        help.Model
}

func initialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Untitled"
	ti.CharLimit = 30
    ti.Width = 30
    ti.TextStyle = li.NewStyle().Bold(true)
    ti.Prompt = "# "
    ti.PromptStyle = li.NewStyle().Bold(true).Faint(true)

	ta := textarea.New()
	ta.Placeholder = "Type a note!"
    ta.SetWidth(40)

	return Model{
		choices: []Choice{
			{Title: "1", Subtitle: ""},
			{Title: "2", Subtitle: ""},
			{Title: "3", Subtitle: ""},
			{Title: "4", Subtitle: ""},
		},
		choiceIndex: 0,
		noteIndex:   -1,
		textinput:   ti,
		textarea:    ta,
		help:        help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	if m.textarea.Focused() {
		return textarea.Blink
	}
	if m.textinput.Focused() {
		return textinput.Blink
	}

	return nil
}

func updateFields(m Model, msg tea.Msg) (Model, []tea.Cmd) {
	cmds := []tea.Cmd{}
	var textinputCmd tea.Cmd
	var textareaCmd tea.Cmd

	m.textinput, textinputCmd = m.textinput.Update(msg)
	m.textarea, textareaCmd = m.textarea.Update(msg)

	cmds = append(cmds, textinputCmd, textareaCmd)
	return m, cmds
}

func updateNote(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, NoteKeyMap.Next):
			if m.textinput.Focused() {
				m.textinput.Blur()
				m.textarea.Focus()
			} else {
				m.textarea.Blur()
				m.textinput.Focus()
			}
		case key.Matches(msg, NoteKeyMap.Cancel):
			m.textinput.Blur()
			m.textarea.Blur()
			m.textinput.SetValue("")
			m.textarea.SetValue("")
			m.noteIndex = -1
		case key.Matches(msg, NoteKeyMap.Save):
			m.noteIndex = -1
			m.textinput.Blur()
			m.textarea.Blur()
			note := &m.choices[m.choiceIndex]
			note.Title = m.textinput.Value()
			note.Subtitle = m.textarea.Value()
			m.textinput.SetValue("")
			m.textarea.SetValue("")
		}
	}
	var fieldCmds []tea.Cmd
	m, fieldCmds = updateFields(m, msg)
	cmds = append(cmds, fieldCmds...)

	return m, tea.Batch(cmds...)
}

func updateList(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ListKeyMap.Next):
			if m.choiceIndex < len(m.choices)-1 {
				m.choiceIndex++
			} else {
				m.choiceIndex = 0
			}

		case key.Matches(msg, ListKeyMap.Prev):
			if m.choiceIndex > 0 {
				m.choiceIndex--
			} else {
				m.choiceIndex = len(m.choices) - 1
			}

		case key.Matches(msg, ListKeyMap.Open):
			m.noteIndex = m.choiceIndex
			m.textinput.SetValue(m.choices[m.noteIndex].Title)
			m.textarea.SetValue(m.choices[m.noteIndex].Subtitle)
			m.textinput.Focus()

		case key.Matches(msg, ListKeyMap.New):
			m.choices = append(m.choices[:m.choiceIndex+1], append([]Choice{{Title: "", Subtitle: ""}}, m.choices[m.choiceIndex+1:]...)...)
			m.choiceIndex++

		case key.Matches(msg, ListKeyMap.Delete):
			if len(m.choices) > 0 {
				m.choices = append(m.choices[:m.choiceIndex], m.choices[m.choiceIndex+1:]...)
			}

		case key.Matches(msg, ListKeyMap.Exit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		if m.viewport.Height == 0 {
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.HighPerformanceRendering = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}
	}

	return m, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.noteIndex < 0 {
		return updateList(m, msg)
	} else {
		return updateNote(m, msg)
	}
}

func viewList(m Model) string {
	str := "\n"

	for i, choice := range m.choices {
		renderedTitle := choice.Title
		if renderedTitle == "" {
			renderedTitle = "Untitled"
		}
		renderedSubtitle := choice.Subtitle
		if renderedSubtitle == "" {
			renderedSubtitle = "Write a note!"
		}

		if m.choiceIndex == i {
			str += listItemSelectedTitle.Render(renderedTitle) + "\n"
			str += listItemSelectedSubtitle.Render(renderedSubtitle) + "\n"
		} else {
			str += listItemTitle.Render(renderedTitle) + "\n"
			str += listItemSubtitle.Render(renderedSubtitle) + "\n"
		}
	}
    str = li.NewStyle().Height(m.viewport.Height).Render(str)
    helpStr := li.NewStyle().Width(m.viewport.Width).Render(m.help.FullHelpView(ListKeyMap.FullHelp()))
    str = li.JoinHorizontal(li.Bottom, str, helpStr)
    return str
}

func viewNote(m Model) string {
	str := ""
	str += m.textinput.View() + "\n\n\n"
    str1_5 := m.textarea.View() + "\n"
    str += str1_5
    helpStr := li.NewStyle().Width(m.viewport.Width).Render(m.help.FullHelpView(NoteKeyMap.FullHelp()))
    str = li.PlaceHorizontal(m.viewport.Width, li.Center, str)
	str = li.PlaceVertical(m.viewport.Height, li.Center, str)
    str = li.JoinVertical(li.Bottom, str, helpStr)
    return str
}

func (m Model) View() string {
	if m.noteIndex < 0 {
		return viewList(m)
	} else {
		return viewNote(m)
	}
}

func alasView(m Model) string {
	str := li.NewStyle().Padding(2, 5).Background(li.Color("#333333")).Foreground(li.Color("#FF0000")).Bold(true).Render("Alas, something has gone amuck!")
	return str
}

func main() {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Print(alasView(m))
		os.Exit(1)
	}
}
