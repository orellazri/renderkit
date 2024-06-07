package main

import (
	"fmt"
	"log"
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
			&cli.PathFlag{
				Name:     "datasource",
				Aliases:  []string{"d"},
				Usage:    "The datasource to use for rendering",
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
			dsFile, err := os.Open(cCtx.String("datasource"))
			if err != nil {
				return fmt.Errorf("open file %s: %s", cCtx.String("datasource"), err)
			}
			defer dsFile.Close()
			ds := datasource.YamlDatasource{}
			data, err := ds.Load(dsFile)
			if err != nil {
				return fmt.Errorf("load from datasource: %s", err)
			}

			outFile, err := os.Create(cCtx.String("output"))
			if err != nil {
				return fmt.Errorf("create file %s: %s", cCtx.String("output"), err)
			}
			defer outFile.Close()

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
