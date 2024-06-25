package chat

import tea "github.com/charmbracelet/bubbletea"

func keyMsgHandler(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, cmd
}
