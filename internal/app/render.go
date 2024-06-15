package app

import (
	"bytes"
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var output io.Writer = os.Stdout
	var err error

	if len(inputString) > 0 { // Render input string
		if len(outputDir) > 0 {
			output, err = createOutputFileWithDir(ctx, filepath.Join(outputDir, "renderkit_output"))
			if err != nil {
				return err
			}
		}
		return a.renderString(inputString, output, data)
	} else if len(inputFile) > 0 { // Render input file
		if len(outputDir) > 0 {
			output, err = createOutputFileWithDir(ctx, filepath.Join(outputDir, filepath.Base(inputFile)))
			if err != nil {
				return err
			}
		}
		return a.renderFile(inputFile, output, excluded, data)
	} else if len(inputDir) > 0 { // Render input directory
		return a.renderDir(inputDir, outputDir, excluded, data)
	}

	return errors.New("unsupported mode")
}

func (a *App) renderDir(inputDirpath string, outputDirpath string, excluded []string, data map[string]any) error {
	err := filepath.WalkDir(inputDirpath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		relPath, err := filepath.Rel(inputDirpath, path)
		if err != nil {
			return fmt.Errorf("get relative path: %s", err)
		}

		var output io.Writer = os.Stdout
		if len(outputDirpath) > 0 {
			output, err = createOutputFileWithDir(ctx, filepath.Join(outputDirpath, relPath))
			if err != nil {
				return err
			}
		}
		if err := a.renderFile(path, output, excluded, data); err != nil {
			return fmt.Errorf("render file %q: %s", path, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("walk directory %q: %s", inputDirpath, err)
	}

	return nil
}

func (a *App) renderFile(inputFilepath string, output io.Writer, excluded []string, data map[string]any) error {
	if slices.Contains(excluded, inputFilepath) {
		return nil
	}

	if err := a.engine.RenderFile(inputFilepath, output, data); err != nil {
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

func createOutputFileWithDir(ctx context.Context, outputFilepath string) (io.Writer, error) {
	outputDirpath := filepath.Dir(outputFilepath)
	if err := os.MkdirAll(outputDirpath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create output directory %s: %s", outputDirpath, err)
	}

	outputFile, err := os.Create(outputFilepath)
	if err != nil {
		return nil, fmt.Errorf("create output file %s: %s", outputFilepath, err)
	}
	go func(ctx context.Context) {
		<-ctx.Done()
		outputFile.Close()
	}(ctx)

	return outputFile, nil
}
