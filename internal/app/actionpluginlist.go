package app

import (
	"fmt"
	"strings"

	"github.com/bcdxn/go-llm/internal/plugins"
	"github.com/bcdxn/go-llm/internal/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

var (
	listStyle = styles.Primary.Margin(1, 2)
	noteStyle = lipgloss.NewStyle().Margin(0)
)

func pluginsList(c *cli.Context) error {
	ps, err := plugins.Find()
	if err != nil {
		return err
	}

	fmt.Println(styles.Title.Render("Installed LLM Plugins:"))

	list := []string{}

	for _, p := range ps {
		list = append(list, fmt.Sprintf("- %s", p.Name))
	}

	fmt.Println(listStyle.Render(strings.Join(list, "\n")))
	fmt.Println(noteStyle.Render(`
You can install more plugins by running:
	go install go-llm-plugin-<plugin-name>
	`))

	fmt.Println(noteStyle.Render(`
You can select an installed plugin to use by running:
	llm plugins select
	`))

	return nil
}
