package app

import (
	"fmt"
	"strings"

	"github.com/orellazri/renderkit/internal/engines"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

type App struct {
	cliApp *cli.App
	engine engines.Engine
}

func NewApp() *App {
	a := App{}

	// Create a list of engine names to display in the CLI help
	engineMapKeys := make([]string, 0, len(enginesMap))
	for k := range enginesMap {
		engineMapKeys = append(engineMapKeys, k)
	}
	enginesListStr := strings.Join(engineMapKeys, ", ")

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Load configuration from YAML file",
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "input",
			Aliases: []string{"i"},
			Usage:   "The input glob to render",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "exclude",
			Aliases:     []string{"x"},
			Usage:       "The glob pattern for files to exclude from rendering",
			DefaultText: "",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "The output directory to write to",
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:    "datasource",
			Aliases: []string{"d"},
			Usage:   "The datasource to use for rendering (scheme://path)",
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:  "data",
			Usage: "The data to use for rendering. Can be used to provide data directly",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "engine",
			Aliases: []string{"e"},
			Usage:   fmt.Sprintf("The templating engine to use for rendering (%s)", enginesListStr),
			Action: func(cCtx *cli.Context, value string) error {
				if _, ok := enginesMap[value]; !ok {
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
	if err := a.validateFlags(
		cCtx.String("input"),
		cCtx.String("output"),
		cCtx.StringSlice("datasource"),
		cCtx.StringSlice("data"),
		cCtx.String("engine"),
	); err != nil {
		if err := cli.ShowAppHelp(cCtx); err != nil {
			return fmt.Errorf("show app help: %s", err)
		}
		return fmt.Errorf("validate flags: %s", err)
	}

	if err := a.setEngine(cCtx.String("engine")); err != nil {
		return fmt.Errorf("set engine: %s", err)
	}

	datasourceUrls, err := a.parseDatasourceUrls(cCtx.StringSlice("datasource"))
	if err != nil {
		return fmt.Errorf("parse datasource URLs: %s", err)
	}

	data, err := a.loadDatasources(datasourceUrls, cCtx.StringSlice("data"), cCtx.Bool("allow-duplicate-keys"))
	if err != nil {
		return fmt.Errorf("load datasources: %s", err)
	}

	inputFiles, err := a.compileGlob(cCtx.String("input"))
	if err != nil {
		return fmt.Errorf("compile input glob: %s", err)
	}

	if len(cCtx.String("exclude")) > 0 {
		excludeFiles, err := a.compileGlob(cCtx.String("exclude"))
		if err != nil {
			return fmt.Errorf("compile exclude glob: %s", err)
		}
		inputFiles = a.excludeFilesFromInput(inputFiles, excludeFiles)
	}

	if err := a.render(
		inputFiles,
		cCtx.String("output"),
		data,
	); err != nil {
		return fmt.Errorf("render: %s", err)
	}

	return nil
}
