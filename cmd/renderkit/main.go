package main

import (
	"fmt"
	"log"
	"os"

	"github.com/orellazri/renderkit/internal/datasource"
	"github.com/orellazri/renderkit/internal/engine"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "renderkit",
		Usage: "A swiss army knife CLI for rendering templates",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"in"},
				Usage:    "The input file to render",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "datasource",
				Aliases:  []string{"ds"},
				Usage:    "The datasource to use for rendering",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"out"},
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
			ds := datasource.NewYamlDatasource()
			data, err := ds.Load(dsFile)
			if err != nil {
				return fmt.Errorf("load from datasource: %s", err)
			}

			outFile, err := os.Create(cCtx.String("output"))
			if err != nil {
				return fmt.Errorf("create file %s: %s", cCtx.String("output"), err)
			}
			defer outFile.Close()
			eng := engine.NewJetEngine()
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
