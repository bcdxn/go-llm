package app

import (
	"context"
	"fmt"

	"github.com/bcdxn/go-llm/internal/chat"
	"github.com/bcdxn/go-llm/internal/config"
	"github.com/bcdxn/go-llm/internal/modelselect"
	"github.com/bcdxn/go-llm/internal/pluginselect"
	"github.com/urfave/cli/v2"
)

func New() *cli.App {
	return &cli.App{
		Version: "1.0.0-rc1",
		Name:    "llm",
		Usage:   "Start an interactive session",
		Before: func(c *cli.Context) error {
			cfg, err := config.Init()
			if err != nil {
				return err
			}

			c.Context = context.WithValue(c.Context, config.CtxConfig{}, cfg)
			return nil
		},
		Action: func(*cli.Context) error {
			_, err := chat.Run()
			return err
		},
		Commands: []*cli.Command{
			{
				Name:  "plugins",
				Usage: "Collection of plugin management commands",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "List the installed plugins",
						Action: func(ctx *cli.Context) error {
							return pluginsList(ctx)
						},
					},
					{
						Name:  "select",
						Usage: "Select a plugin from a list of your installed plugins",
						Action: func(ctx *cli.Context) error {
							_, err := pluginselect.Run(ctx)
							if err != nil {
								return err
							}

							cfg, err := config.Load()
							ctx.Context = context.WithValue(ctx.Context, config.CtxConfig{}, cfg)
							// if the model was reset we can go ahead and prompt the user
							if cfg.DefaultModel == "" {
								_, err = modelselect.Run(ctx)
							}

							return err
						},
					},
				},
			},
			{
				Name:  "models",
				Usage: "Collection of model management commands",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "List the models supported by the selected plugin",
						Action: func(ctx *cli.Context) error {
							fmt.Println("coming soon...")
							return nil
						},
					},
					{
						Name:  "select",
						Usage: "Select a model to use",
						Action: func(ctx *cli.Context) error {
							_, err := modelselect.Run(ctx)
							return err
						},
					},
				},
			},
			{
				Name:  "chat",
				Usage: "Start an interactive session with a model",
				Action: func(*cli.Context) error {
					_, err := chat.Run()
					return err
				},
			},
			{
				Name:    "message",
				Aliases: []string{"msg"},
				Usage:   "Send a message in-line to a model",
				Action: func(*cli.Context) error {
					fmt.Println("comming soon...")
					return nil
				},
			},
		},
	}
}
