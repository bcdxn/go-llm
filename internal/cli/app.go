package cli

import (
	"fmt"
	"plugin"

	"github.com/bcdxn/go-llm/internal/chat"
	"github.com/urfave/cli/v2"
)

func New() *cli.App {
	return &cli.App{
		Version: "1.0.0-rc1",
		Name:    "llm",
		Usage:   "Start an interactive session",
		Action: func(*cli.Context) error {
			_, err := chat.Run()
			return err
		},
		Commands: []*cli.Command{
			{
				Name: "list",
				Subcommands: []*cli.Command{
					{
						Name:  "plugins",
						Usage: "List the installed plugins",
						Action: func(*cli.Context) error {
							fmt.Println("coming soon...")
							return nil
						},
					},
					{
						Name:  "models",
						Usage: "List the installed plugins",
						Action: func(*cli.Context) error {
							fmt.Println("coming soon...")
							return nil
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
			{
				Name:  "lp",
				Usage: "load plugin test",
				Action: func(*cli.Context) error {
					p, err := plugin.Open("/Users/bdxn/Documents/2024-interview-prep/projects/go-llm-plugin-openai/go-llm-plugin-openai.so")
					if err != nil {
						fmt.Println("something bad happened loading the plugin...")
						fmt.Println(err)
						panic(err)
					}

					v, err := p.Lookup("V")
					if err != nil {
						panic(err)
					}
					f, err := p.Lookup("F")
					if err != nil {
						panic(err)
					}

					*v.(*int) = 7
					f.(func())()

					return nil
				},
			},
		},
	}
}
