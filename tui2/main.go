package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/StefanSchroeder/bubbletd"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pborman/ansi"
)

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2).Align(lipgloss.Left)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	nonactionColorHi    = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"}
	nonactionColorLo    = lipgloss.AdaptiveColor{Light: "#770000", Dark: "#770000"}

	actionColorHi    = lipgloss.AdaptiveColor{Light: "#00FF00", Dark: "#00FF00"}
	actionColorLo    = lipgloss.AdaptiveColor{Light: "#007700", Dark: "#007700"}

	inactiveTabStyle  =     lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(  activeTabBorder, true).Foreground(lipgloss.Color("111"))

	redactiveTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true).Foreground(nonactionColorHi)
	redinactiveTabStyle    = inactiveTabStyle.Copy().Border(inactiveTabBorder, true).Foreground(nonactionColorLo)

	blueactiveTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true).Foreground(actionColorHi)
	blueinactiveTabStyle    = inactiveTabStyle.Copy().Border(inactiveTabBorder, true).Foreground(actionColorLo)

	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 2).Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop()

	nonactionStyleA     = lipgloss.NewStyle().Foreground(lipgloss.Color("111"))
	nonactionStyleB     = focusedStyle.Copy()
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
	activeTab  int
	TextInputs []textinput.Model
	table      table.Model
	textarea   textarea.Model
	indexstore int
	btd        bubbletd.Bubbletd
}

var isNonactionable = map[string]bool{
    "Trash": true,
    "Reference":   true,
    "Later": true,
}
var isActionable = map[string]bool{
    "Quick": true,
    "Float":   true,
    "Dairy": true,
    "Someone": true,
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

	if m.focusIndex == 1 || m.focusIndex == 0 {
		// Changing selected entry. Retrieve entry for textarea
		current_table_row2 := m.table.SelectedRow()
		if len(current_table_row2) > 0 {
			current_table_index2, _ := strconv.Atoi(current_table_row2[0])
			tf := m.btd[current_table_index2].Desc
			m.textarea.SetValue(tf)
		} else {
			m.textarea.SetValue("")
		}
	}

	return tea.Batch(cmds...)
}

func (m *model) build_table(a []string, gotocursor int, filter_state string) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Task", Width: 40},
	}

	rows := []table.Row{}

	for i := len(a) - 1; i >= 0; i-- {
		j := a[i]
		if m.btd[i].State == filter_state {
			rows = append(rows, []string{fmt.Sprint(i), j})
		}
	}

	tb := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(10),
	)
	if gotocursor != -1 {
		tb.SetCursor(gotocursor)
	}

	current_table_row2 := m.table.SelectedRow()
	if len(current_table_row2) > 0 {
		current_table_index2, _ := strconv.Atoi(current_table_row2[0])
		tf := m.btd[current_table_index2].Desc
		m.textarea.SetValue(tf)
	} else {
		m.textarea.SetValue("")
	}

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

// CleanInputString removed ANSI control characters and the
// prompt that for some reason is part of the input.
func CleanInputString(s string) string {
	entered_text := strings.TrimSpace(s)
	entered_text = strings.SplitN(entered_text, " ", 2)[1]
	d := []byte(entered_text)
	d2, _ := ansi.Strip(d)
	return (string(d2))
}

func (m *model) MoveEntryToState(s string) {
	if m.focusIndex == 1 {
		current_table_row := m.table.SelectedRow()
		if len(current_table_row) > 0 {
			current_table_index, _ := strconv.Atoi(current_table_row[0])

			m.btd[current_table_index].State = s

			titles := m.btd.GetTitles()
			m.table = m.build_table(titles, m.table.Cursor(), m.Tabs[m.activeTab])
			m.View()
			m.table.View()
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+k":
			m.TextInputs[0].SetValue("")
			return m, nil
		case "ctrl+f":
			m.MoveEntryToState("Float")
			return m, nil
		case "ctrl+q":
			m.MoveEntryToState("Quick")
			return m, nil
		case "ctrl+a":
			m.MoveEntryToState("Attic")
			return m, nil
		case "ctrl+s":
			m.MoveEntryToState("Someone")
			return m, nil
		case "ctrl+d":
			m.MoveEntryToState("Dairy")
			return m, nil
		case "ctrl+r":
			m.MoveEntryToState("Reference")
			return m, nil
		case "ctrl+l":
			m.MoveEntryToState("Later")
			return m, nil
		case "ctrl+t":
			m.MoveEntryToState("Trash")
			return m, nil
		case "ctrl+c", "esc":
			m.btd.WriteConfig()
			return m, tea.Quit
		case "f1":
			m.activeTab = max(m.activeTab-1, 0)
			titles := m.btd.GetTitles()
			m.table = m.build_table(titles, m.table.Cursor(), m.Tabs[m.activeTab])
			m.updateInputs(msg)
			return m, nil
		case "f2":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			titles := m.btd.GetTitles()
			m.table = m.build_table(titles, m.table.Cursor(), m.Tabs[m.activeTab])
			m.updateInputs(msg)
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

				// Leaving textarea. Storing entry
				current_table_row := m.table.SelectedRow()
				if len(current_table_row) > 0 {
					s := fmt.Sprint(m.textarea.Value())
					m.btd.SetDesc("desc " + current_table_row[0] + " " + s)
				}
			}

			return m, nil
		case "enter":
			if m.focusIndex == 0 {
				entered_text := CleanInputString(m.TextInputs[0].View())

				if m.indexstore == -1 {
					// This is a new entry
					m.btd.AddTask("add " + entered_text)
					titles := m.btd.GetTitles()

					m.activeTab = 0
					m.table = m.build_table(titles, m.table.Cursor(), "Inbox")
					m.updateInputs(msg)

					m.table.SetCursor(0)
					m.textarea.SetValue("")
				} else {
					// This is a rewrite entry
					m.btd.EditTitle("edit " + fmt.Sprint(m.indexstore) + " " + entered_text)
					titles := m.btd.GetTitles()
					m.table = m.build_table(titles, m.table.Cursor(), "Inbox")
				}
				m.TextInputs[0].SetValue("")

				// clear flag
				m.indexstore = -1
			}
			if m.focusIndex == 1 {
				// get text, fill into inputfield, raise flag that this is a correction
				current_table_row := m.table.SelectedRow()
				if len(current_table_row) > 0 {
					current_table_index, _ := strconv.Atoi(current_table_row[0])
					current_table_string := current_table_row[1]
					m.TextInputs[0].SetValue(current_table_string)

					m.indexstore = current_table_index
					// when this is a correction, we are not going to create a new entry, but use the registered flag
					// move focus to textentry.
					m.focusIndex = 0
					m.TextInputs[0].Focus()
					m.table.Blur()
					m.textarea.Blur()
				}
			}
		}
	}

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

		// The inefficieny of this makes my skin crawl.
		cnt := 0
		for _, j := range m.btd {
			if t == j.State {
				cnt += 1
			}
		}
		cntS := fmt.Sprintf(" (%02d)", cnt)
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle.Copy()
			if isNonactionable[t] {
				style = redactiveTabStyle.Copy()
			}
			if isActionable[t] {
				style = blueactiveTabStyle.Copy()
			}
		} else {
			style = inactiveTabStyle.Copy()
			if isNonactionable[t] {
				style = redinactiveTabStyle.Copy()
			}
			if isActionable[t] {
				style = blueinactiveTabStyle.Copy()
			}
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
		renderedTabs = append(renderedTabs, style.Render(t+cntS))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	if m.activeTab == 0 || m.activeTab > 0 {
		x := ""
		if m.indexstore != -1 {
			x += fmt.Sprint("rewriting (", m.indexstore, ") ")
		}
		x += fmt.Sprint(m.TextInputs[0].View())
		x += "\n"
		x += "\n"
		x += m.table.View()
		x += "\n"
		x += "\n"
		x += "Description\n"
		x += "\n"
		x += m.textarea.View()

		doc.WriteString(windowStyle.Width(4 + (lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(x))

	} else if m.activeTab == 1 {
		doc.WriteString(windowStyle.Width(4 + (lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render("abc"))
	} else if m.activeTab == 2 {
		doc.WriteString(windowStyle.Width(4 + (lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render("abc"))
	} else {
		doc.WriteString(windowStyle.Width(4 + (lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render("abc"))
	}

	if m.focusIndex == 0 {
		doc.WriteString("\nEnter task")
	}
	/*if m.focusIndex == 1 {
		doc.WriteString("\nIs this actionable?")
	}
	if m.focusIndex == 2 {
		sr := fmt.Sprintf("%v", m.table.SelectedRow())
		doc.WriteString("\nDesc for " + sr)
	}*/

	return docStyle.Render(doc.String())
}

func main() {

	btd := bubbletd.New()
	btd.ReadConfig()

	btd.Review()

fmt.Println(btd)

	tia := textarea.New()
	tia.Placeholder = "Elaboration of task..."

	tabs := []string{"Inbox", "Trash", "Reference", "Later", "Quick", "Float", "Dairy", "Someone", "Attic"}
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	m := model{
		Tabs:       tabs,
		TextInputs: make([]textinput.Model, 1),
		table:      table.New(),
		textarea:   tia,
		indexstore: -1,
		btd:        *btd}

	m.table = m.build_table(btd.GetTitles(), 0, "Inbox")

	var t textinput.Model
	for i := range m.TextInputs {
		t = textinput.New()
		t.CharLimit = 32

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

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	fmt.Println("Good-bye")
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
