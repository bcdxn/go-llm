package main

import (
	"fmt"

	"github.com/bcdxn/go-llm/internal/plugins"
)

// import (
// 	"log"
// 	"os"

// 	"github.com/bcdxn/go-llm/internal/cli"
// )

func main() {
	fs, _ := plugins.Find()

	for _, f := range fs {
		_, err := plugins.Load(f)
		if err != nil {
			panic(err)
		}
		fmt.Println("calling F")
		// p.F()
	}
	// app := cli.New()

	// if err := app.Run(os.Args); err != nil {
	// 	log.Fatal(err)
	// }
}
