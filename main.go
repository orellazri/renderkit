package main

import (
	"log"
	"os"

	"github.com/orellazri/renderkit/internal/app"
)

func main() {
	app := app.NewApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
