package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

func (a *App) render(
	inputDir string,
	inputFile string,
	outputDir string,
	excluded []string,
	data map[string]any,
) error {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("create output directory %s: %s", outputDir, err)
	}

	if len(inputDir) > 0 {
		return a.renderDir(inputDir, outputDir, excluded, data)
	} else if len(inputFile) > 0 {
		return a.renderFile(inputFile, filepath.Join(outputDir, filepath.Base(inputFile)), excluded, data)
	}

	return errors.New("unsupported mode")
}

func (a *App) renderDir(inputDirname, outputDirname string, excluded []string, data map[string]any) error {
	err := filepath.WalkDir(inputDirname, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(inputDirname, path)
		if err != nil {
			return fmt.Errorf("get relative path: %s", err)
		}

		// Create output subdirectory
		err = os.MkdirAll(filepath.Join(outputDirname, filepath.Dir(relPath)), os.ModePerm)
		if err != nil {
			return fmt.Errorf("create output directory %q: %s", outputDirname, err)
		}

		outFilename := filepath.Join(outputDirname, relPath)
		if err := a.renderFile(path, outFilename, excluded, data); err != nil {
			return fmt.Errorf("render file %q: %s", path, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("walk directory %q: %s", inputDirname, err)
	}

	return nil
}

func (a *App) renderFile(inFilename, outFilename string, excluded []string, data map[string]any) error {
	if slices.Contains(excluded, inFilename) {
		return nil
	}

	outFile, err := os.Create(outFilename)
	if err != nil {
		return fmt.Errorf("create output file %s: %s", outFilename, err)
	}
	defer outFile.Close()

	if err := a.engine.RenderFile(inFilename, outFile, data); err != nil {
		return fmt.Errorf("render template: %s", err)
	}

	return nil
}
