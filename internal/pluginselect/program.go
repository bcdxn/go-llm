// This file represents a fully encapsulated bubbletea program that allows a user to select a plugin
// from a list of available plugins

package pluginselect

import (
	"github.com/bcdxn/go-llm/internal/config"
	"github.com/bcdxn/go-llm/internal/plugins"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

/* Program Entry Point
------------------------------------------------------------------------------------------------- */

func Run(ctx *cli.Context) (tea.Model, error) {
	return tea.NewProgram(getInitialModel(ctx), tea.WithAltScreen()).Run()
}

/* Model
------------------------------------------------------------------------------------------------- */

type model struct {
	list         list.Model
	prevSelected plugins.PluginListItem
	selected     plugins.PluginListItem
	width        int // window width
	height       int // window height
	cfg          config.Config
}

/* Component
------------------------------------------------------------------------------------------------- */

func (m model) Init() tea.Cmd {
	return getPluginsCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case getPluginsMsg:
		return getPluginsHandler(m)
	case tea.KeyMsg:
		return keyMsgHandler(m, msg)
	case tea.WindowSizeMsg:
		return windowSizeMsgHandler(m, msg)
	case updateConfigMsg:
		return updateConfigMsgHandler(m)
	default:
		var cmd tea.Cmd
		return m, cmd
	}
}

func (m model) View() string {
	return "\n" + m.list.View()
}

func getInitialModel(ctx *cli.Context) model {
	items := []list.Item{}

	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.SetSpacing(0)
	l := list.New(items, d, 40, 10)
	l.Title = "Select a plugin to use:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	cfg, ok := ctx.Context.Value(config.CtxConfig{}).(config.Config)
	if !ok {
		cfg = config.Config{}
	}

	return model{
		prevSelected: cfg.DefaultPlugin,
		list:         l,
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

func updateConfigCmd() tea.Cmd {
	return func() tea.Msg {
		return updateConfigMsg{}
	}
}

/* Handlers
------------------------------------------------------------------------------------------------- */

func getPluginsHandler(m model) (model, tea.Cmd) {
	ps, _ := plugins.Find()

	items := []list.Item{}

	for _, p := range ps {
		items = append(items, item(p))
	}

	m.list.SetItems(items)

	return m, nil
}

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
	m.cfg.DefaultPlugin = plugins.PluginListItem(m.selected)
	if m.cfg.DefaultPlugin.Path != m.prevSelected.Path {
		m.cfg.DefaultModel = "" // reset the default model when plugin changes
	}
	config.Persist(m.cfg)
	return m, tea.Quit
}

/* Messages
------------------------------------------------------------------------------------------------- */

type getPluginsMsg struct{}
type updateConfigMsg struct{}
