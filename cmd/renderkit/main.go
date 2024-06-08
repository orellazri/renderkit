package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/orellazri/renderkit/internal/datasource"
	"github.com/orellazri/renderkit/internal/engine"
	"github.com/urfave/cli/v2"
)

func main() {
	// Create a list of engine names to display in the CLI help
	engineMapKeys := make([]string, 0, len(engine.EnginesMap))
	for k := range engine.EnginesMap {
		engineMapKeys = append(engineMapKeys, k)
	}
	enginesListStr := strings.Join(engineMapKeys, ", ")

	app := &cli.App{
		Name:  "renderkit",
		Usage: "A swiss army knife CLI tool for rendering templates",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "The input file to render",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     "datasource",
				Aliases:  []string{"d"},
				Usage:    "The datasource to use for rendering (scheme://path)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "engine",
				Aliases:  []string{"e"},
				Usage:    fmt.Sprintf("The engine to use for rendering (%s)", enginesListStr),
				Required: true,
				Action: func(cCtx *cli.Context, value string) error {
					if _, ok := engine.EnginesMap[value]; !ok {
						return fmt.Errorf("engine %s not supported. supported engines:  %s", value, enginesListStr)
					}
					return nil
				},
			},
			&cli.PathFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "The output file to write to",
				Required: true,
			},
		},
		Action: func(cCtx *cli.Context) error {
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
