package chat

import tea "github.com/charmbracelet/bubbletea"

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return keyMsgHandler(m, msg)
	default:
		return m, nil
	}
}

func (m model) View() string {
	str := "I'm a bubble tea component"
	return str
}

func Run() (tea.Model, error) {
	return tea.NewProgram(model{}, tea.WithAltScreen()).Run()
}
