package app

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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
	var output io.Writer
	if len(outputDir) > 0 {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return fmt.Errorf("create output directory %s: %s", outputDir, err)
		}
	} else {
		output = os.Stdout
	}

	if len(inputString) > 0 { // Render input string
		if len(outputDir) > 0 {
			outputFilename := filepath.Join(outputDir, "renderkit_output")
			outputFile, err := os.Create(outputFilename)
			if err != nil {
				return fmt.Errorf("create output file %s: %s", outputFilename, err)
			}
			defer outputFile.Close()
			output = outputFile
		}

		return a.renderString(inputString, output, data)
	} else if len(inputFile) > 0 { // Render input file
		if len(outputDir) > 0 {
			outputFilename := filepath.Join(outputDir, filepath.Base(inputFile))
			outputFile, err := os.Create(outputFilename)
			if err != nil {
				return fmt.Errorf("create output file %s: %s", outputFilename, err)
			}
			defer outputFile.Close()
			output = outputFile
		}

		return a.renderFile(inputFile, output, excluded, data)
	} else if len(inputDir) > 0 { // Render input directory
		return a.renderDir(inputDir, outputDir, excluded, data)
	}

	return errors.New("unsupported mode")
}

func (a *App) renderDir(inputDirname string, outputDirname string, excluded []string, data map[string]any) error {
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

		var output io.Writer
		if len(outputDirname) > 0 {
			// Create output subdirectory
			err = os.MkdirAll(filepath.Join(outputDirname, filepath.Dir(relPath)), os.ModePerm)
			if err != nil {
				return fmt.Errorf("create output directory %q: %s", outputDirname, err)
			}

			outputFilename := filepath.Join(outputDirname, relPath)
			outputFile, err := os.Create(outputFilename)
			if err != nil {
				return fmt.Errorf("create output file %q: %s", outputFilename, err)
			}
			output = outputFile
		} else {
			output = os.Stdout
		}
		if err := a.renderFile(path, output, excluded, data); err != nil {
			return fmt.Errorf("render file %q: %s", path, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("walk directory %q: %s", inputDirname, err)
	}

	return nil
}

func (a *App) renderFile(inputFilename string, output io.Writer, excluded []string, data map[string]any) error {
	if slices.Contains(excluded, inputFilename) {
		return nil
	}

	if err := a.engine.RenderFile(inputFilename, output, data); err != nil {
		return fmt.Errorf("render template: %s", err)
	}

	return nil
}

func (a *App) renderString(inputString string, output io.Writer, data map[string]any) error {
	if err := a.engine.Render(bytes.NewReader([]byte(inputString)), output, data); err != nil {
		return fmt.Errorf("render template: %s", err)
	}

	return nil
}
