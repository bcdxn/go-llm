package styles

import "github.com/charmbracelet/lipgloss"

var (
	Title = lipgloss.NewStyle().Foreground(
		lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1)
	Accent  = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	Primary = lipgloss.NewStyle()
)
