package app

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	llm "github.com/bcdxn/go-llm/internal"
	"github.com/bcdxn/go-llm/internal/shared"
	"github.com/hashicorp/go-plugin"
	"github.com/urfave/cli/v2"
)

func modelsList(ctx *cli.Context) error {
	var (
		cfg       = llm.MustGetConfigFromContext(ctx.Context)
		l         = shared.MustGetLoggerFromContext(ctx.Context, "modelslist")
		pluginMap = map[string]plugin.Plugin{
			"llm": &shared.LLMPlugin{},
		}
		client = plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: shared.DefaultHandshakeConfig,
			Plugins:         pluginMap,
			Logger:          l,
			Cmd:             exec.Command(cfg.DefaultPlugin.Path),
		})
	)
	defer client.Kill()

	l.Trace("creating RPC Client")
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	l.Trace("successfully created RPC Client")
	l.Trace("requesting Plugin from RPC Client")

	raw, err := rpcClient.Dispense("llm")
	if err != nil {
		log.Fatal(err)
	}
	l.Trace("successfully requested Plugin from RPC Client")

	llm := raw.(shared.LLM)

	selectedPlugin := cfg.DefaultPlugin.Name
	ms := llm.GetModels(cfg.Plugins[selectedPlugin])
	l.Debug("Successfully fetched models", "models", ms)

	list := []string{}

	for _, m := range ms {
		list = append(list, fmt.Sprintf("- %s", m))
	}

	fmt.Println(listStyle.Render(strings.Join(list, "\n")))

	fmt.Println(noteStyle.Render(`
You can select a supported model to use by running:
	llm models select
	`))

	return nil
}
