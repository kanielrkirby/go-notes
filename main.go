package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var c = map[string]lipgloss.Color{
	"primary":          lipgloss.Color("#F4442E"),
	"primary-light":    lipgloss.Color("#FC9E4F"),
	"primary-lighter":  lipgloss.Color("#EDD382"),
	"primary-lightest": lipgloss.Color("#F2F3AE"),
	"black":            lipgloss.Color("#020122"),
	"white":            lipgloss.Color("#FEFEF7"),
}

var (
	listItemTitle            = lipgloss.NewStyle().MarginLeft(1).PaddingLeft(1).Bold(true).Foreground(c["white"]).Border(lipgloss.HiddenBorder(), false, false, false, true)
	listItemSubtitle         = lipgloss.NewStyle().MarginLeft(1).MarginBottom(1).PaddingLeft(1).Foreground(c["white"]).Faint(true).Border(lipgloss.HiddenBorder(), false, false, false, true)
	listItemSelectedTitle    = lipgloss.NewStyle().MarginLeft(1).PaddingLeft(1).Bold(true).Foreground(c["primary-light"]).Border(lipgloss.ThickBorder(), false, false, false, true).BorderForeground(c["primary-light"])
	listItemSelectedSubtitle = lipgloss.NewStyle().MarginLeft(1).MarginBottom(1).PaddingLeft(1).Foreground(c["primary-light"]).Faint(true).Border(lipgloss.ThickBorder(), false, false, false, true).BorderForeground(c["primary-light"])
)

type Choice struct {
	Title    string
	Subtitle string
}

type model struct {
	viewport    viewport.Model
	textinput   textinput.Model
	textarea    textarea.Model
	choices     []Choice
	choiceIndex int
	noteIndex   int
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Untitled"
	ti.CharLimit = 30

	ta := textarea.New()
	ta.Placeholder = "Type a note!"

	return model{
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
	}
}

func (m model) Init() tea.Cmd {
	if m.textarea.Focused() {
		return textarea.Blink
	}
	if m.textinput.Focused() {
		return textinput.Blink
	}

	return nil
}

func viewList(m model) string {
	str := "\n"

	for i, choice := range m.choices {
		title := choice.Title
		if title == "" {
			title = "Untitled"
		}
		subtitle := choice.Subtitle
		if subtitle == "" {
			subtitle = "Write a note!"
		}

		if m.choiceIndex == i {
			str += listItemSelectedTitle.Render(title) + "\n"
			str += listItemSelectedSubtitle.Render(subtitle) + "\n"
		} else {
			str += listItemTitle.Render(title) + "\n"
			str += listItemSubtitle.Render(subtitle) + "\n"
		}
	}
	return str
}

func updateFields(m model, msg tea.Msg) (model, []tea.Cmd) {
	cmds := []tea.Cmd{}
	var textinputCmd tea.Cmd
	var textareaCmd tea.Cmd

	m.textinput, textinputCmd = m.textinput.Update(msg)
	m.textarea, textareaCmd = m.textarea.Update(msg)

	cmds = append(cmds, textinputCmd, textareaCmd)
	return m, cmds
}

func updateNote(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type.String() {
		case tea.KeyTab.String():
			if m.textinput.Focused() {
				m.textinput.Blur()
				m.textarea.Focus()
			} else {
				m.textarea.Blur()
				m.textinput.Focus()
			}
		case tea.KeyCtrlC.String():
			return m, tea.Quit
		case tea.KeyEnter.String():
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

func updateList(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit

		case tea.KeyUp.String(), "k", tea.KeyShiftTab.String():
			if m.choiceIndex > 0 {
				m.choiceIndex--
			} else {
				m.choiceIndex = len(m.choices) - 1
			}

		case tea.KeyShiftTab.String(), "j", tea.KeyTab.String():
			if m.choiceIndex < len(m.choices)-1 {
				m.choiceIndex++
			} else {
				m.choiceIndex = 0
			}

		case tea.KeyEnter.String(), tea.KeySpace.String(), "l":
			m.noteIndex = m.choiceIndex
			m.textinput.SetValue(m.choices[m.noteIndex].Title)
			m.textarea.SetValue(m.choices[m.noteIndex].Subtitle)
			m.textinput.Focus()

		case "n":
            m.choices = append(m.choices[:m.choiceIndex+1], append([]Choice{{Title: "", Subtitle: ""}}, m.choices[m.choiceIndex+1:]...)...)
			m.choiceIndex++

		case "d":
			if len(m.choices) > 0 {
				m.choices = append(m.choices[:m.choiceIndex], m.choices[m.choiceIndex+1:]...)
			}
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.noteIndex < 0 {
		return updateList(m, msg)
	} else {
		return updateNote(m, msg)
	}
}

func viewNote(m model) string {
	str := ""
	str += m.textinput.View() + "\n"
	str += m.textarea.View() + "\n"
	return lipgloss.PlaceVertical(m.viewport.Height, lipgloss.Center, lipgloss.PlaceHorizontal(m.viewport.Width, lipgloss.Center, str))
}

func (m model) View() string {
	if m.noteIndex < 0 {
		return viewList(m)
	} else {
		return viewNote(m)
	}
}

func alasView(m model) string {
	str := lipgloss.NewStyle().Padding(2, 5).Background(lipgloss.Color("#333333")).Foreground(lipgloss.Color("#FF0000")).Bold(true).Render("Alas, something has gone amuck!")
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
