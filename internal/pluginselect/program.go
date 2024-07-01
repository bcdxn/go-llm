// This file represents a fully encapsulated bubbletea program that allows a user to select a plugin
// from a list of available plugins

package pluginselect

import (
	"fmt"

	llm "github.com/bcdxn/go-llm/internal"
	"github.com/bcdxn/go-llm/internal/plugins"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

/* Program Entry Point
------------------------------------------------------------------------------------------------- */

func Run(c *cli.Context) (tea.Model, error) {
	return tea.NewProgram(getInitialModel(c), tea.WithAltScreen()).Run()
}

/* Model
------------------------------------------------------------------------------------------------- */

type model struct {
	l            llm.Logger
	list         list.Model
	prevSelected plugins.PluginListItem
	selected     plugins.PluginListItem
	width        int // window width
	height       int // window height
	cfg          llm.Config
}

/* Component
------------------------------------------------------------------------------------------------- */

func (m model) Init() tea.Cmd {
	m.l.Debug("pluginselect Init")
	return getPluginsCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.l.Debug("message received", "msg", fmt.Sprintf("%T", msg))
	switch msg := msg.(type) {
	case getPluginsMsg:
		return getPluginsHandler(m)
	case tea.KeyMsg:
		return keyMsgHandler(m, msg)
	case tea.WindowSizeMsg:
		return windowSizeMsgHandler(m, msg)
	case selectedMsg:
		return selectedMsgHandler(m)
	default:
		var cmd tea.Cmd
		return m, cmd
	}
}

func (m model) View() string {
	return "\n" + m.list.View()
}

func getInitialModel(c *cli.Context) model {
	var (
		l     = llm.MustGetLoggerFromContext(c.Context, "pluginselect")
		items = []list.Item{}
		d     = list.NewDefaultDelegate()
	)

	d.ShowDescription = false
	d.SetSpacing(0)
	list := list.New(items, d, 40, 10)
	list.Title = "Select a plugin to use:"
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)

	cfg := llm.MustGetConfigFromContext(c.Context)

	return model{
		l:            l,
		prevSelected: cfg.DefaultPlugin,
		list:         list,
		cfg:          cfg,
	}
}

type item plugins.PluginListItem

func (i item) FilterValue() string { return "" }
func (i item) Title() string       { return i.Name }
func (i item) Description() string { return i.Path }

/* Commands
------------------------------------------------------------------------------------------------- */

func getPluginsCmd() tea.Cmd {
	return func() tea.Msg {
		return getPluginsMsg{}
	}
}

func selectedCmd() tea.Cmd {
	return func() tea.Msg {
		return selectedMsg{}
	}
}

/* Handlers
------------------------------------------------------------------------------------------------- */

func getPluginsHandler(m model) (model, tea.Cmd) {
	m.l.Debug("plugins handler")
	ps, _ := plugins.Find()
	m.l.Debug("plugins", "plugins", ps)

	items := []list.Item{}

	for _, p := range ps {
		items = append(items, item(p))
	}

	m.list.SetItems(items)

	return m, nil
}

// Handle keystrokes to navigate the list or quit the app.
func keyMsgHandler(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	m.l.Debug("keyMsgHandler", "key", msg.String())
	switch keypress := msg.String(); keypress {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "enter":
		i, ok := m.list.SelectedItem().(item)
		if ok {
			m.selected = plugins.PluginListItem(i)
		}
		return m, selectedCmd()
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// Handle window resize events and update the size of the list accordingly.
func windowSizeMsgHandler(m model, msg tea.WindowSizeMsg) (model, tea.Cmd) {
	m.l.Debug("windowSizeMsgHandler", "width", msg.Width, "height", msg.Height)
	h, v := lipgloss.NewStyle().GetFrameSize()
	m.width = msg.Width - h
	m.height = msg.Height - v

	m.list.SetWidth(m.width)
	m.list.SetHeight(m.height)

	return m, nil
}

// Handle configuration update events by writing the selected plugin from the list to the persistent
// config file.
func selectedMsgHandler(m model) (model, tea.Cmd) {
	m.l.Debug("selectedMsgHandler", "selected", m.selected)
	m.cfg.DefaultPlugin = plugins.PluginListItem(m.selected)
	if m.cfg.DefaultPlugin.Path != m.prevSelected.Path {
		m.cfg.DefaultModel = "" // reset the default model when plugin changes
	}
	llm.PersistConfig(m.cfg)
	return m, tea.Quit
}

/* Messages
------------------------------------------------------------------------------------------------- */

type getPluginsMsg struct{}
type selectedMsg struct{}
