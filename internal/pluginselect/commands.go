package pluginselect

import tea "github.com/charmbracelet/bubbletea"

func updateConfigCmd() tea.Cmd {
	return func() tea.Msg {
		return updateConfigMsg{}
	}
}
