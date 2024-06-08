package app

import (
	"fmt"
	"os"
	"strings"

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
	if err := validateFlags(cCtx); err != nil {
		return fmt.Errorf("validate flags: %s", err)
	}

	datasourceUrls, err := parseDatasourceUrls(cCtx)
	if err != nil {
		return fmt.Errorf("parse datasource URLs: %s", err)
	}

	data, err := loadDatasources(datasourceUrls)
	if err != nil {
		return fmt.Errorf("create datasources: %s", err)
	}

	outFile, err := os.Create(cCtx.String("output"))
	if err != nil {
		return fmt.Errorf("create output file %s: %s", cCtx.String("output"), err)
	}
	defer outFile.Close()

	if err := renderTemplate(cCtx, outFile, data); err != nil {
		return fmt.Errorf("render template: %s", err)
	}

	return nil
}
