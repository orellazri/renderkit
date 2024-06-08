package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/orellazri/renderkit/internal/engine"
	"github.com/urfave/cli/v2"
)

func (a *App) render(cCtx *cli.Context, data map[string]any) error {
	switch a.mode {
	case ModeFileToFile:
		return a.renderFileToFile(cCtx, data)
	case ModeFilesToDir:
		return a.renderFilesToDir(cCtx, data)
	case ModeDirToDir:
		return a.renderDirToDir(cCtx, data)
	}

	return nil
}

func (a *App) renderFileToFile(cCtx *cli.Context, data map[string]any) error {
	outFile, err := os.Create(cCtx.String("output"))
	if err != nil {
		return fmt.Errorf("create output file %s: %s", cCtx.String("output"), err)
	}
	defer outFile.Close()

	if err := a.renderTemplate(cCtx, cCtx.StringSlice("input")[0], outFile, data); err != nil {
		return fmt.Errorf("render template: %s", err)
	}

	return nil
}

func (a *App) renderFilesToDir(cCtx *cli.Context, data map[string]any) error {
	if err := os.MkdirAll(cCtx.String("output-dir"), os.ModePerm); err != nil {
		return fmt.Errorf("create output directory %s: %s", cCtx.String("output-dir"), err)
	}

	for _, inFilename := range cCtx.StringSlice("input") {
		outFilename := filepath.Join(cCtx.String("output-dir"), filepath.Base(inFilename))
		outFile, err := os.Create(outFilename)
		if err != nil {
			return fmt.Errorf("create output file %s: %s", outFilename, err)
		}
		defer outFile.Close()
		if err := a.renderTemplate(cCtx, inFilename, outFile, data); err != nil {
			return fmt.Errorf("render template: %s", err)
		}
	}

	return nil
}

func (a *App) renderDirToDir(cCtx *cli.Context, data map[string]any) error {
	if err := os.MkdirAll(cCtx.String("output-dir"), os.ModePerm); err != nil {
		return fmt.Errorf("create output directory %s: %s", cCtx.String("output-dir"), err)
	}

	entries, err := os.ReadDir(cCtx.String("input-dir"))
	if err != nil {
		return fmt.Errorf("read input directory %s: %s", cCtx.String("input-dir"), err)
	}

	for _, inFilename := range entries {
		outFilename := filepath.Join(cCtx.String("output-dir"), filepath.Base(inFilename.Name()))
		outFile, err := os.Create(outFilename)
		if err != nil {
			return fmt.Errorf("create output file %s: %s", outFilename, err)
		}
		defer outFile.Close()
		if err := a.renderTemplate(cCtx, filepath.Join(cCtx.String("input-dir"), inFilename.Name()), outFile, data); err != nil {
			return fmt.Errorf("render template: %s", err)
		}
	}

	return nil
}

func (a *App) renderTemplate(cCtx *cli.Context, inFilename string, outFile *os.File, data map[string]any) error {
	eng, ok := engine.EnginesMap[cCtx.String("engine")]
	if !ok {
		return fmt.Errorf("engine %s not found", cCtx.String("engine"))
	}
	if err := eng.Render(inFilename, outFile, data); err != nil {
		return fmt.Errorf("render: %s", err)
	}

	return nil
}
