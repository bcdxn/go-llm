package chat

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	llm "github.com/bcdxn/go-llm/internal"
	"github.com/bcdxn/go-llm/internal/shared"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-plugin"
	"github.com/urfave/cli/v2"
)

/* Program Entry Point
--------------------------------------------------------------------------------------------------*/

func Run(ctx *cli.Context) (tea.Model, error) {
	return tea.NewProgram(getInitialModel(ctx), tea.WithAltScreen()).Run()
}

/* Model
--------------------------------------------------------------------------------------------------*/

type model struct {
	l              shared.Logger
	cfg            llm.Config
	viewport       viewport.Model
	textarea       textarea.Model
	conversationId string
	messages       []string
	model          string
	width          int // window width
	height         int // window height
}

/* Component Behavior
--------------------------------------------------------------------------------------------------*/

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return windowSizeMsgHandler(m, msg)
	case tea.KeyMsg:
		return keyMsgHandler(m, msg)
	case sendPromptMsg:
		return sendPromptHandler(m, msg)
	default:
		return m, nil
	}
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func getInitialModel(c *cli.Context) model {
	var (
		l   = shared.MustGetLoggerFromContext(c.Context, "modelselect")
		cfg = llm.MustGetConfigFromContext(c.Context)
		ta  = textarea.New()
		vp  = viewport.New(30, 5)
	)

	ta.Placeholder = "Send a message..."
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)
	ta.SetWidth(30)
	ta.SetHeight(5)
	ta.Focus()

	vp.SetContent("")

	return model{
		l:        l,
		cfg:      cfg,
		textarea: ta,
		viewport: vp,
	}
}

/* Commands
--------------------------------------------------------------------------------------------------*/

func sendPromptCmd(msg string) tea.Cmd {
	return func() tea.Msg {
		return sendPromptMsg{msg}
	}
}

/* Handlers
--------------------------------------------------------------------------------------------------*/

func keyMsgHandler(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	var taCmd tea.Cmd
	m.textarea, taCmd = m.textarea.Update(msg)

	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit
	case tea.KeyEnter:
		prompt := m.textarea.Value()
		// m.messages = append(m.messages, prompt)
		m.textarea.Reset()
		return m, sendPromptCmd(prompt)
	}
	return m, taCmd
}

func windowSizeMsgHandler(m model, msg tea.WindowSizeMsg) (model, tea.Cmd) {
	h, v := lipgloss.NewStyle().GetFrameSize()
	m.width = msg.Width - h
	m.height = msg.Height - v

	m.textarea.SetWidth(m.width)
	m.textarea.SetHeight(5)

	m.viewport.Width = m.width
	m.viewport.Height = m.height - 8

	return m, nil
}

func sendPromptHandler(m model, msg sendPromptMsg) (model, tea.Cmd) {
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

	llmp := raw.(shared.LLM)
	ms := llmp.SendMessages(shared.MessageParam{
		Config: m.cfg.Plugins[m.cfg.DefaultPlugin.Name],
		Model:  m.cfg.DefaultModel,
		Messages: []shared.Message{
			{Role: "User", Content: msg.prompt},
		},
	})

	for _, res := range ms {
		m.messages = append(m.messages, res.Content)
	}

	m.viewport.SetContent(strings.Join(m.messages, "\n"))

	return m, nil
}

/* Messages
--------------------------------------------------------------------------------------------------*/

type sendPromptMsg struct {
	prompt string
}
