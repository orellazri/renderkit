package app

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
			Usage:   "Template string to render",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "input-file",
			Aliases: []string{"f"},
			Usage:   "Template input file to render",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "input-dir",
			Aliases: []string{"d"},
			Usage:   "Template input directory to render",
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:        "exclude",
			Aliases:     []string{"x"},
			Usage:       "Exclude files/directories using path-based glob patterns",
			DefaultText: "",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Output directory to write to",
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:    "datasource",
			Aliases: []string{"ds"},
			Usage:   "Datasource to use for rendering (scheme://path)",
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:  "data",
			Usage: "Data to use for rendering. Can be used to provide data directly",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "engine",
			Aliases: []string{"e"},
			Usage:   fmt.Sprintf("Templating engine to use for rendering (%s)", enginesListStr),
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
	var inputString string

	// Read from stdin into an input string; and if empty, from input flag
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		var stdinBytes []byte
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			stdinBytes = append(stdinBytes, scanner.Bytes()...)
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Failed to read from stdin: %s", err)
		}
		inputString = string(stdinBytes)
	} else if len(cCtx.String("input")) > 0 {
		inputString = cCtx.String("input")
	}

	if err := a.validateFlags(
		inputString,
		cCtx.String("input-dir"),
		cCtx.String("input-file"),
		cCtx.StringSlice("datasource"),
		cCtx.StringSlice("data"),
		cCtx.StringSlice("exclude"),
		cCtx.String("engine"),
	); err != nil {
		if err := cli.ShowAppHelp(cCtx); err != nil {
			return fmt.Errorf("show app help: %s", err)
		}
		return fmt.Errorf("validate flags: %s", err)
	}

	if eng, ok := enginesMap[cCtx.String("engine")]; !ok {
		a.engine = enginesMap["gotemplates"]
	} else {
		a.engine = eng
	}

	datasourceUrls, err := a.parseDatasourceUrls(cCtx.StringSlice("datasource"))
	if err != nil {
		return fmt.Errorf("parse datasource URLs: %s", err)
	}

	data, err := a.loadDatasources(datasourceUrls, cCtx.StringSlice("data"), cCtx.Bool("allow-duplicate-keys"))
	if err != nil {
		return fmt.Errorf("load datasources: %s", err)
	}

	excludePaths := []string{}
	excludeFileGlobs := []string{}
	if len(cCtx.StringSlice("exclude")) > 0 {
		excludePaths, excludeFileGlobs, err = a.aggregateExcludePatterns(cCtx.StringSlice("exclude"))
		if err != nil {
			return fmt.Errorf("aggregate exclude patterns: %s", err)
		}
	}

	if err := a.render(
		inputString,
		cCtx.String("input-dir"),
		cCtx.String("input-file"),
		cCtx.String("output"),
		excludePaths,
		excludeFileGlobs,
		data,
	); err != nil {
		return fmt.Errorf("render: %s", err)
	}

	return nil
}
