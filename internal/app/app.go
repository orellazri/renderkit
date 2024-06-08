package app

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/orellazri/renderkit/internal/datasource"
	"github.com/orellazri/renderkit/internal/engine"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func NewApp() *cli.App {
	// Create a list of engine names to display in the CLI help
	engineMapKeys := make([]string, 0, len(engine.EnginesMap))
	for k := range engine.EnginesMap {
		engineMapKeys = append(engineMapKeys, k)
	}
	enginesListStr := strings.Join(engineMapKeys, ", ")

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Usage: "Load configuration from YAML file",
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "input",
			Aliases: []string{"i"},
			Usage:   "The input file to render",
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:    "datasource",
			Aliases: []string{"d"},
			Usage:   "The datasource to use for rendering (scheme://path)",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "engine",
			Aliases: []string{"e"},
			Usage:   fmt.Sprintf("The engine to use for rendering (%s)", enginesListStr),
			Action: func(cCtx *cli.Context, value string) error {
				if _, ok := engine.EnginesMap[value]; !ok {
					return fmt.Errorf("engine %s is not supported. supported engines: %s", value, enginesListStr)
				}
				return nil
			},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "The output file to write to",
		}),
	}

	app := &cli.App{
		Name:   "renderkit",
		Usage:  "A swiss army knife CLI tool for rendering templates",
		Flags:  flags,
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config")),
		Action: run,
	}

	return app
}

func run(cCtx *cli.Context) error {
	// Validate flags
	if err := validateFlags(cCtx); err != nil {
		return fmt.Errorf("validate flags: %s", err)
	}

	// Parse datasource URLs
	datasourceUrls := make([]*url.URL, len(cCtx.StringSlice("datasource")))
	for i, ds := range cCtx.StringSlice("datasource") {
		url, err := url.Parse(ds)
		if err != nil {
			return fmt.Errorf("parse datasource %s: %s", ds, err)
		}
		datasourceUrls[i] = url
	}

	// Create datasources
	data := make(map[string]any)
	for _, url := range datasourceUrls {
		ds, err := datasource.CreateDatasourceFromURL(url)
		if err != nil {
			return fmt.Errorf("create datasource %s: %s", url, err)
		}

		dsData, err := ds.Load()
		if err != nil {
			return fmt.Errorf("load datasource %s: %s", url, err)
		}

		// Merge with data dictionary
		for k, v := range dsData {
			data[k] = v
		}
	}

	// Create output file
	outFile, err := os.Create(cCtx.String("output"))
	if err != nil {
		return fmt.Errorf("create file %s: %s", cCtx.String("output"), err)
	}
	defer outFile.Close()

	// Render the template
	eng, ok := engine.EnginesMap[cCtx.String("engine")]
	if !ok {
		return fmt.Errorf("engine %s not found", cCtx.String("engine"))
	}
	if err := eng.Render(cCtx.String("input"), outFile, data); err != nil {
		return fmt.Errorf("render: %s", err)
	}

	return nil
}

func validateFlags(cCtx *cli.Context) error {
	if len(cCtx.String("input")) == 0 {
		return fmt.Errorf("input file is required")
	}

	if len(cCtx.StringSlice("datasource")) == 0 {
		return fmt.Errorf("datasource is required")
	}

	if len(cCtx.String("engine")) == 0 {
		return fmt.Errorf("engine is required")
	}

	if len(cCtx.String("output")) == 0 {
		return fmt.Errorf("output file is required")
	}

	return nil
}
