package pluginselect

import (
	"github.com/bcdxn/go-llm/internal/config"
	"github.com/bcdxn/go-llm/internal/plugins"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Handle keystrokes to navigate the list or quit the app.
func keyMsgHandler(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "enter":
		i, ok := m.list.SelectedItem().(item)
		if ok {
			m.selected = plugins.PluginListItem(i)
		}
		return m, updateConfigCmd()
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// Handle window resize events and update the size of the list accordingly.
func windowSizeMsgHandler(m model, msg tea.WindowSizeMsg) (model, tea.Cmd) {
	h, v := lipgloss.NewStyle().GetFrameSize()
	m.width = msg.Width - h
	m.height = msg.Height - v

	m.list.SetWidth(m.width)
	m.list.SetHeight(m.height)

	return m, nil
}

// Handle configuration update events by writing the selected plugin from the list to the persistent
// config file.
func updateConfigMsgHandler(m model) (model, tea.Cmd) {
	m.cfg.SelectedPlugin = plugins.PluginListItem(m.selected)
	config.Persist(m.cfg)
	return m, tea.Quit
}
