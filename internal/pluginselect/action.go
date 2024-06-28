package pluginselect

import (
	"context"

	"github.com/bcdxn/go-llm/internal/config"
	"github.com/bcdxn/go-llm/internal/modelselect"
	"github.com/urfave/cli/v2"
)

func Action(ctx *cli.Context) error {
	_, err := Run(ctx)
	if err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// if the model was reset we can go ahead and prompt the user for the model
	if cfg.DefaultModel == "" {
		ctx.Context = context.WithValue(ctx.Context, config.CtxConfig{}, cfg)
		_, err := modelselect.Run(ctx)
		if err != nil {
			return err
		}
	}

	cfg, err = config.Load()
	if err != nil {
		return err
	}

	return err
}
