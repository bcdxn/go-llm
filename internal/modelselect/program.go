package modelselect

import (
	"fmt"
	"os"
	"os/exec"

	llm "github.com/bcdxn/go-llm/internal"
	"github.com/bcdxn/go-llm/internal/shared"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-plugin"
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
	l        llm.Logger
	width    int // window width
	height   int // window height
	list     list.Model
	cfg      llm.Config
	selected string
}

/* Component
------------------------------------------------------------------------------------------------- */

func (m model) Init() tea.Cmd {
	return fetchModelsCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.l.Debug("message received", "msg", fmt.Sprintf("%T", msg))
	var cmds []tea.Cmd

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case fetchModelsMsg:
		return fetchModelsHandler(m)
	case tea.KeyMsg:
		return keyMsgHandler(m, msg)
	case tea.WindowSizeMsg:
		return windowSizeMsgHandler(m, msg)
	case selectedMsg:
		return selectedMsgHandler(m)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return "\n" + m.list.View()
}

func getInitialModel(c *cli.Context) model {
	var (
		l     = llm.MustGetLoggerFromContext(c.Context, "modelselect")
		d     = list.NewDefaultDelegate()
		items = []list.Item{}
		cfg   = llm.MustGetConfigFromContext(c.Context)
	)
	defer l.Close()

	d.ShowDescription = false
	d.SetSpacing(0)
	list := list.New(items, d, 30, 10)
	list.Title = fmt.Sprintf("Select a model to use from the %s plugin:", cfg.DefaultPlugin.Name)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)

	return model{
		l:    l,
		list: list,
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

func selectedCmd() tea.Cmd {
	return func() tea.Msg {
		return selectedMsg{}
	}
}

/* Handlers
------------------------------------------------------------------------------------------------- */

func fetchModelsHandler(m model) (model, tea.Cmd) {
	var (
		pluginMap = map[string]plugin.Plugin{
			"llm": &shared.LLMPlugin{},
		}
	)

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.DefaultHandshakeConfig,
		Plugins:         pluginMap,
		Logger:          m.l,
		Cmd:             exec.Command(m.cfg.DefaultPlugin.Path),
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		m.l.Error("error initializing RPC client", err)
		os.Exit(1)
	}

	raw, err := rpcClient.Dispense("llm")
	if err != nil {
		m.l.Error("error initializing plugin", err)
		os.Exit(1)
	}

	llm := raw.(shared.LLM)
	ms := llm.GetModels()
	m.l.Debug("successfully retrieved models", "models", ms)

	models := []list.Item{}

	for _, m := range ms {
		models = append(models, item(m))
	}

	m.l.Trace("updating list model with models")
	m.list.SetItems(models)

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
			m.selected = string(i)
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

// Handle configuration update events by writing the selected model from the list to the persistent
// config file.
func selectedMsgHandler(m model) (model, tea.Cmd) {
	m.l.Debug("selectedCmdMsgHandler", "selected", m.selected)
	m.cfg.DefaultModel = m.selected
	llm.PersistConfig(m.cfg)
	return m, tea.Quit
}

/* Messages
------------------------------------------------------------------------------------------------- */

type selectedMsg struct{}
type fetchModelsMsg struct{}
