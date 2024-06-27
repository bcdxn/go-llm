package main

import (
	"log"
	"os"

	"github.com/bcdxn/go-llm/internal/app"
)

func main() {
	app := app.New()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
