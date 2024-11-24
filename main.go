package main

import (
	"log"
	"os"
	"runtime/debug"

	"github.com/orellazri/renderkit/internal/app"
)

// This version variable is set at compile time using ldflags
var version = "dev"

func main() {
	// If the version is not set at compile time, try to get it from the build info
	info, ok := debug.ReadBuildInfo()
	if ok && info.Main.Version != "(devel)" {
		version = info.Main.Version
	}

	app := app.NewApp(version)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
