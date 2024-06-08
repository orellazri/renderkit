package app

import (
	"fmt"
	"strings"

	"github.com/orellazri/renderkit/internal/engine"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

type Mode int

const (
	ModeFileToFile Mode = iota
	ModeFilesToDir
	ModeDirToDir
)

type App struct {
	cliApp *cli.App
	mode   Mode
	engine engine.Engine
}

func NewApp() *App {
	a := App{}

	// Create a list of engine names to display in the CLI help
	engineMapKeys := make([]string, 0, len(engine.EnginesMap))
	for k := range engine.EnginesMap {
		engineMapKeys = append(engineMapKeys, k)
	}
	enginesListStr := strings.Join(engineMapKeys, ", ")

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Load configuration from YAML file",
		},
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:    "input",
			Aliases: []string{"i"},
			Usage:   "The input file to render",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "The output file to write to",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "input-dir",
			Usage: "The input directory to render",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "output-dir",
			Usage: "The output directory to write to",
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
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "allow-duplicate-keys",
			Usage:       "Allow duplicate keys in datasources. If set, the last value found will be used",
			DefaultText: "false",
		}),
	}

	app := &cli.App{
		Name:   "renderkit",
		Usage:  "A swiss army knife CLI tool for rendering templates",
		Flags:  flags,
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config")),
		Action: a.run,
	}

	a.cliApp = app
	return &a
}

func (a *App) Run(args []string) error {
	return a.cliApp.Run(args)
}

func (a *App) run(cCtx *cli.Context) error {
	if err := a.validateFlags(cCtx); err != nil {
		return fmt.Errorf("validate flags: %s", err)
	}

	if err := a.setMode(cCtx.StringSlice("input"), cCtx.String("input-dir")); err != nil {
		return fmt.Errorf("set mode: %s", err)
	}

	if err := a.setEngine(cCtx.String("engine")); err != nil {
		return fmt.Errorf("set engine: %s", err)
	}

	datasourceUrls, err := a.parseDatasourceUrls(cCtx.StringSlice("datasource"))
	if err != nil {
		return fmt.Errorf("parse datasource URLs: %s", err)
	}

	data, err := a.loadDatasources(datasourceUrls, cCtx.Bool("allow-duplicate-keys"))
	if err != nil {
		return fmt.Errorf("load datasources: %s", err)
	}

	if err := a.render(
		cCtx.StringSlice("input"),
		cCtx.String("output"),
		cCtx.String("input-dir"),
		cCtx.String("output-dir"),
		data,
	); err != nil {
		return fmt.Errorf("render: %s", err)
	}

	return nil
}
