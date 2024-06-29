package app

import (
	"fmt"
	"strings"

	"github.com/bcdxn/go-llm/internal/logger"
	"github.com/bcdxn/go-llm/internal/plugins"
	"github.com/bcdxn/go-llm/internal/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

var (
	listStyle = styles.Primary.Margin(1, 2)
	noteStyle = lipgloss.NewStyle().Margin(0)
)

func pluginsList(ctx *cli.Context) error {
	l, ok := ctx.Context.Value(logger.CtxLogger{}).(*logger.Logger)
	if !ok {
		logger.SimpleLogFatal("unable to fetch logger from context")
	}
	ll := l.Named("pluginslist")

	ll.Trace("finding plugis")
	ps, err := plugins.Find()
	if err != nil {
		return err
	}
	ll.Debug("foud plugins", "plugins", ps)

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
