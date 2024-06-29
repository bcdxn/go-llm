package app

import (
	"fmt"

	"github.com/bcdxn/go-llm/internal/logger"
	"github.com/urfave/cli/v2"
)

func modelsList(ctx *cli.Context) error {
	l, ok := ctx.Context.Value(logger.CtxLogger{}).(*logger.Logger)
	if !ok {
		logger.SimpleLogFatal("unable to fetch logger from context")
	}
	ll := l.Named("modelslist")

	ll.Trace("Fetching models")
	fmt.Println("coming soon...")
	ll.Debug("Fetched models")

	return nil
}
