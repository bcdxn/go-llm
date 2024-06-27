package pluginselect

import (
	"fmt"
	"io"
	"strings"

	tealogger "github.com/bcdxn/go-llm/internal"
	"github.com/bcdxn/go-llm/internal/plugins"
	"github.com/bcdxn/go-llm/internal/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	logger            = tealogger.New(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = styles.Accent.PaddingLeft(2)
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return keyMsgHandler(m, msg)
	case tea.WindowSizeMsg:
		return windowSizeMsgHandler(m, msg)
	default:
		var cmd tea.Cmd
		return m, cmd
	}
}

func (m model) View() string {
	return "\n" + m.list.View()
}

func Run() (tea.Model, error) {
	return tea.NewProgram(getInitialModel(), tea.WithAltScreen()).Run()
}

func getInitialModel() model {
	ps, err := plugins.Find()
	if err != nil {
		logger.LogFatal(err)
	}

	items := []list.Item{}

	for _, p := range ps {
		items = append(items, item(p))
	}

	l := list.New(items, itemDelegate{}, 40, 10)
	l.Title = "Select a plugin to use"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return model{
		plugins: ps,
		list:    l,
	}
}

type item plugins.PluginListItem

func (i item) FilterValue() string { return "" }
func (i item) Title() string       { return i.Name }
func (i item) Description() string { return i.Path }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
