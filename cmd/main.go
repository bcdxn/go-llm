package main

// import (
// 	"fmt"

// 	"github.com/bcdxn/go-llm/internal/logger"
// 	"github.com/hashicorp/go-hclog"
// )

import (
	"log"
	"os"

	"github.com/bcdxn/go-llm/internal/app"
)

func main() {
	// l, err := logger.New(hclog.Trace)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// l.Debug("hello logging world")
	// l.Info("hello logging world")
	// l.Error("an error")

	app := app.New()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
