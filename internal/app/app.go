package app

import (
	"fmt"
	"os"
	"path/filepath"
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
			Name:  "input-dir",
			Usage: "The input files directory to render",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "The output file to write to",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "output-dir",
			Usage: "The output files directory to render",
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

	data, err := loadDatasources(datasourceUrls, cCtx.Bool("allow-duplicate-keys"))
	if err != nil {
		return fmt.Errorf("load datasources: %s", err)
	}

	if err := handleInputs(cCtx, data); err != nil {
		return fmt.Errorf("handle inputs: %s", err)
	}

	return nil
}

func handleInputs(cCtx *cli.Context, data map[string]any) error {
	if cCtx.StringSlice("input") != nil {
		return handleFiles(cCtx, data)
	} else if cCtx.String("input-dir") != "" {
		return handleDirectories(cCtx, data)
	}
	return nil
}

func handleDirectories(cCtx *cli.Context, data map[string]any) error {
	inputDirPathname := cCtx.String("input-dir")
	outputDirPathname := cCtx.String("output-dir")

	inputDir, err := os.Open(inputDirPathname)
	if err != nil {
		return fmt.Errorf("open input directory %s: %s", inputDirPathname, err)
	}
	defer inputDir.Close()
	inputDirFileNames, err := inputDir.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("read input directory %s: %s", inputDirPathname, err)
	}

	if err := os.MkdirAll(outputDirPathname, 0755); err != nil {
		return fmt.Errorf("create output directory %s: %s", outputDirPathname, err)
	}

	for _, f := range inputDirFileNames {
		inputFilePathname := fmt.Sprintf("%s/%s", inputDirPathname, f)
		inputFile, err := os.Open(inputFilePathname)
		if err != nil {
			return fmt.Errorf("open input file %s: %s", inputFilePathname, err)
		}
		defer inputFile.Close()
		outputFilePath := fmt.Sprintf("%s/%s", outputDirPathname, f)
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			return fmt.Errorf("create output file %s: %s", outputFilePath, err)
		}
		defer outputFile.Close()
		if err := renderTemplate(cCtx, inputFile, outputFile, data); err != nil {
			return fmt.Errorf("render template: %s", err)
		}
	}
	return nil
}

func handleFiles(cCtx *cli.Context, data map[string]any) error {
	var outputFilePathname string
	for _, f := range cCtx.StringSlice("input") {
		inputFile, err := os.Open(f)
		if err != nil {
			return fmt.Errorf("open input file %s: %s", f, err)
		}
		defer inputFile.Close()
		if len(cCtx.StringSlice("input")) == 1 {
			outputFilePathname = cCtx.String("output")
		} else {
			outputFilePathname = fmt.Sprintf("%s/%s", cCtx.String("output-dir"), filepath.Base(f))
		}
		if err := os.MkdirAll(cCtx.String("output-dir"), 0755); err != nil {
			return fmt.Errorf("create output directory %s: %s", cCtx.String("output-dir"), err)
		}
		outputFile, err := os.Create(outputFilePathname)
		if err != nil {
			return fmt.Errorf("create output file %s: %s", outputFilePathname, err)
		}
		defer outputFile.Close()
		if err := renderTemplate(cCtx, inputFile, outputFile, data); err != nil {
			return fmt.Errorf("render template: %s", err)
		}
	}
	return nil
}
