package pluginselect

import (
	llm "github.com/bcdxn/go-llm/internal"
	"github.com/bcdxn/go-llm/internal/modelselect"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	_, err := Run(c)
	if err != nil {
		return err
	}

	cfg, err := llm.LoadConfig()
	if err != nil {
		return err
	}

	// if the model was reset we can go ahead and prompt the user for the model
	if cfg.DefaultModel == "" {
		c.Context = llm.SetConfigInContext(c.Context, cfg)
		_, err := modelselect.Run(c)
		if err != nil {
			return err
		}
	}

	cfg, err = llm.LoadConfig()
	if err != nil {
		return err
	}

	return err
}
