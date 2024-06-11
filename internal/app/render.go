package app

import (
	"fmt"
	"os"
	"path/filepath"
)

func (a *App) render(input []string, output, inputDir, outputDir string, data map[string]any) error {
	switch a.mode {
	case ModeFileToFile:
		return a.renderFileToFile(input[0], output, data)
	case ModeFilesToDir:
		return a.renderFilesToDir(input, outputDir, data)
	case ModeDirToDir:
		return a.renderDirToDir(inputDir, outputDir, data)
	}

	return nil
}

func (a *App) renderFileToFile(input string, output string, data map[string]any) error {
	return a.renderFile(input, output, data)
}

func (a *App) renderFilesToDir(input []string, outputDir string, data map[string]any) error {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("create output directory %s: %s", outputDir, err)
	}

	for _, inFilename := range input {
		outFilename := filepath.Join(outputDir, filepath.Base(inFilename))
		if err := a.renderFile(inFilename, outFilename, data); err != nil {
			return fmt.Errorf("render file: %s", err)
		}
	}

	return nil
}

func (a *App) renderDirToDir(inputDir, outputDir string, data map[string]any) error {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("create output directory %s: %s", outputDir, err)
	}

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("read input directory %s: %s", inputDir, err)
	}

	for _, entry := range entries {
		inFilename := filepath.Join(inputDir, entry.Name())
		outFilename := filepath.Join(outputDir, entry.Name())
		if err := a.renderFile(inFilename, outFilename, data); err != nil {
			return fmt.Errorf("render file: %s", err)
		}
	}

	return nil
}

func (a *App) renderFile(inFilename, outFilename string, data map[string]any) error {
	outFile, err := os.Create(outFilename)
	if err != nil {
		return fmt.Errorf("create output file %s: %s", outFilename, err)
	}
	defer outFile.Close()

	if err := a.engine.Render(inFilename, outFile, data); err != nil {
		return fmt.Errorf("render template: %s", err)
	}

	return nil
}
