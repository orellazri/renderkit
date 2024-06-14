package app

import (
	"fmt"
	"os"
	"path/filepath"
)

func (a *App) render(inputFiles []string, outputDir string, data map[string]any) error {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("create output directory %s: %s", outputDir, err)
	}

	for _, inputFilename := range inputFiles {
		outFilename := filepath.Join(outputDir, filepath.Base(inputFilename))
		if err := a.renderFile(inputFilename, outFilename, data); err != nil {
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

	if err := a.engine.RenderFile(inFilename, outFile, data); err != nil {
		return fmt.Errorf("render template: %s", err)
	}

	return nil
}
