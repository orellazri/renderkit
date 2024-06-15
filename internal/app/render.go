package app

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

func (a *App) render(
	inputString string,
	inputDir string,
	inputFile string,
	outputDir string,
	excluded []string,
	data map[string]any,
) error {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("create output directory %s: %s", outputDir, err)
	}

	if len(inputString) > 0 {
		return a.renderString(inputString, filepath.Join(outputDir, "renderkit_output"), data)
	} else if len(inputDir) > 0 {
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

		outputFilename := filepath.Join(outputDirname, relPath)
		if err := a.renderFile(path, outputFilename, excluded, data); err != nil {
			return fmt.Errorf("render file %q: %s", path, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("walk directory %q: %s", inputDirname, err)
	}

	return nil
}

func (a *App) renderFile(inputFilename, outputFilename string, excluded []string, data map[string]any) error {
	if slices.Contains(excluded, inputFilename) {
		return nil
	}

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("create output file %s: %s", outputFilename, err)
	}
	defer outputFile.Close()

	if err := a.engine.RenderFile(inputFilename, outputFile, data); err != nil {
		return fmt.Errorf("render template: %s", err)
	}

	return nil
}

func (a *App) renderString(inputString, outputFilename string, data map[string]any) error {
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("create output file %s: %s", outputFilename, err)
	}
	defer outputFile.Close()

	if err := a.engine.Render(bytes.NewReader([]byte(inputString)), outputFile, data); err != nil {
		return fmt.Errorf("render template: %s", err)
	}

	return nil
}
