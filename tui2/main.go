package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2).Align(lipgloss.Left)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 2).Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop()

	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

type model struct {
	focusIndex int
	Tabs       []string
	TabContent []string
	activeTab  int
	TextInputs []textinput.Model
	NewTask    string
	table      table.Model
	textarea   textarea.Model
	data       []string
}

func (m model) Init() tea.Cmd {

	return nil
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.TextInputs))

	// Only text inputs with Focus set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.TextInputs {
		m.TextInputs[i], cmds[i] = m.TextInputs[i].Update(msg)
	}

	m.table, _ = m.table.Update(msg)
	m.textarea, _ = m.textarea.Update(msg)
	return tea.Batch(cmds...)
}

func build_table(a []string) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Task", Width: 40},
	}

	rows := []table.Row{}

	for i := len(a) - 1; i >= 0; i-- {
		j := a[i]
		rows = append(rows, []string{fmt.Sprint(i), j})
	}

	/*	for i, j := range a {
		rows = append( rows, []string{fmt.Sprint(i), j } )
	}*/

	tb := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(5),
	)

	/*s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	tb.SetStyles(s)*/
	return tb
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "f2":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "tab":

			m.focusIndex++

			if m.focusIndex == 1 {
				m.TextInputs[0].Blur()
				m.table.Focus()
				m.textarea.Blur()
			} else if m.focusIndex == 2 {
				m.TextInputs[0].Blur()
				m.table.Blur()
				m.textarea.Focus()
			} else if m.focusIndex == 3 {
				m.focusIndex = 0
				m.TextInputs[0].Focus()
				m.table.Blur()
				m.textarea.Blur()
			}

			return m, nil
		case "f1":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		case "enter":
			entered_text := strings.TrimSpace(m.TextInputs[0].View())
			entered_text = strings.SplitN(entered_text, " ", 2)[1]

			m.data = append(m.data, entered_text)
			m.table = build_table(m.data)
			m.TextInputs[0].SetValue("")
		}

	}
	/*if m.focusIndex > len(m.TextInputs) {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = len(m.TextInputs)
	}*/

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	cmds := make([]tea.Cmd, len(m.TextInputs))
	for i := 0; i <= len(m.TextInputs)-1; i++ {
		if i == m.focusIndex {
			// Set focused state
			cmds[i] = m.TextInputs[i].Focus()
			m.TextInputs[i].PromptStyle = focusedStyle
			m.TextInputs[i].TextStyle = focusedStyle
			continue
		}
		m.TextInputs[i].Blur()
		m.TextInputs[i].PromptStyle = noStyle
		m.TextInputs[i].TextStyle = noStyle
	}

	// m.TextInput, _ = m.TextInput.Update(msg)
	return m, cmd
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func (m model) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	if m.activeTab == 0 {
		x := fmt.Sprint(m.TextInputs[0].View())
		x += "\n"
		x += "\n"
		x += m.table.View()
		x += "\n"
		x += "\n"
		x += "Description\n"
		x += "\n"
		x += m.textarea.View()

		doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(x))

		//doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.TextInputs[0].View()))
	} else if m.activeTab == 1 {
		doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.TabContent[m.activeTab]))
	} else if m.activeTab == 2 {
		doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.TabContent[m.activeTab]))
	} else {
		doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.TabContent[m.activeTab]))
	}
	return docStyle.Render(doc.String())
}

func main() {

	tia := textarea.New()
	tia.Placeholder = "Once upon a time..."

	tb := build_table([]string{})

	tabs := []string{"Inbox    ", "Trash    ", "Reference     ", "Deferred", "Quick", "Queue", "Calendar", "Delegated"}
	tabContent := []string{"inbox", "Trash", "Reference", "Deferred", "Quick", "Queue", "Cal", "Del"}
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	m := model{Tabs: tabs, TabContent: tabContent, TextInputs: make([]textinput.Model, 1), table: tb, textarea: tia, data: []string{}}

	var t textinput.Model
	for i := range m.TextInputs {
		t = textinput.New()
		t.CharLimit = 22

		switch i {
		case 0:
			t.Placeholder = "New Task"
			t.Focus()
		case 1:
			t.Placeholder = "Email"
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Password"
			t.EchoCharacter = '•'
		}

		m.TextInputs[i] = t
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
