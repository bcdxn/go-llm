package main

import (
	"log"
	"os"

	"github.com/bcdxn/go-llm/internal/cli"
)

func main() {
	app := cli.New()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
