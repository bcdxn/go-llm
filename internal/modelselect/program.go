package modelselect

import (
	"fmt"

	"github.com/bcdxn/go-llm/internal/config"
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
	width    int // window width
	height   int // window height
	list     list.Model
	cfg      config.Config
	selected string
}

/* Component
------------------------------------------------------------------------------------------------- */

func (m model) Init() tea.Cmd {
	return fetchModelsCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	f, _ := tea.LogToFile("test.log", "")
	defer f.Close()
	f.WriteString(fmt.Sprintf("wat: %T - %v\n", msg, msg))
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case fetchModelsMsg:
		return fetchModelsHandler(m)
	case tea.KeyMsg:
		return keyMsgHandler(m, msg)
	case tea.WindowSizeMsg:
		return windowSizeMsgHandler(m, msg)
	case updateConfigMsg:
		return updateConfigMsgHandler(m)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return "\n" + m.list.View()
}

func getInitialModel(ctx *cli.Context) model {
	cfg, ok := ctx.Context.Value(config.CtxConfig{}).(config.Config)
	if !ok {
		cfg = config.Config{}
	}

	items := []list.Item{}

	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.SetSpacing(0)
	l := list.New(items, d, 30, 10)
	l.Title = fmt.Sprintf("Select a model to use from %s:", cfg.DefaultPlugin.Name)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return model{
		list: l,
		cfg:  cfg,
	}
}

type item string

func (i item) FilterValue() string { return "" }
func (i item) Title() string       { return string(i) }
func (i item) Description() string { return string(i) }

/* Commands
------------------------------------------------------------------------------------------------- */

func fetchModelsCmd() tea.Cmd {
	return func() tea.Msg {
		return fetchModelsMsg{}
	}
}

func updateConfigCmd() tea.Cmd {
	return func() tea.Msg {
		return updateConfigMsg{}
	}
}

/* Handlers
------------------------------------------------------------------------------------------------- */

func fetchModelsHandler(m model) (model, tea.Cmd) {
	// todo: make plugin call to fetch models
	models := []list.Item{
		item("gpt-3.5-turbo-1106"),
		item("gpt-3.5-turbo-16k"),
		item("gpt-3.5-turbo-0125"),
		item("gpt-3.5-turbo"),
		item("gpt-3.5-turbo-instruct"),
		item("gpt-3.5-turbo-instruct-0914"),
	}

	m.list.SetItems(models)

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
			m.selected = string(i)
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

// Handle configuration update events by writing the selected model from the list to the persistent
// config file.
func updateConfigMsgHandler(m model) (model, tea.Cmd) {
	m.cfg.DefaultModel = m.selected
	config.Persist(m.cfg)
	return m, tea.Quit
}

/* Messages
------------------------------------------------------------------------------------------------- */

type updateConfigMsg struct{}
type fetchModelsMsg struct{}
