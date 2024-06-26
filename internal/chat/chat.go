package chat

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return windowSizeMsgHandler(m, msg)
	case tea.KeyMsg:
		return keyMsgHandler(m, msg)
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

func Run() (tea.Model, error) {
	return tea.NewProgram(getInitialModel(), tea.WithAltScreen()).Run()
}

func getInitialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)
	ta.SetWidth(30)
	ta.SetHeight(5)
	ta.Focus()

	vp := viewport.New(30, 5)
	vp.SetContent("")

	return model{
		textarea: ta,
		viewport: vp,
	}
}
