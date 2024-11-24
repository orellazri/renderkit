package main

import (
	"log"
	"os"

	"github.com/orellazri/renderkit/internal/app"
)

// This version variable is set at compile time using ldflags
var version = "dev"

func main() {
	app := app.NewApp(version)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
