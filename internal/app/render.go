package app

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	"github.com/gobwas/glob"
)

func (a *App) render(
	inputString string,
	inputDir string,
	inputFile string,
	outputDir string,
	excludedPaths []string,
	excludeSingularGlobs []string,
	data map[string]any,
) error {
	var output io.Writer = os.Stdout
	var closer func()
	var err error

	if len(inputString) > 0 { // Render input string
		if len(outputDir) > 0 {
			output, closer, err = createOutputFileWithDir(filepath.Join(outputDir, "renderkit_output"))
			if err != nil {
				return err
			}
			defer closer()
		}
		return a.renderString(inputString, output, data)
	} else if len(inputFile) > 0 { // Render input file
		if len(outputDir) > 0 {
			output, closer, err = createOutputFileWithDir(filepath.Join(outputDir, filepath.Base(inputFile)))
			if err != nil {
				return err
			}
			defer closer()
		}
		return a.renderFile(inputFile, output, data)
	} else if len(inputDir) > 0 { // Render input directory
		// if outputDir is empty
		return a.renderDir(inputDir, outputDir, excludedPaths, excludeSingularGlobs, data)
	}

	return errors.New("unsupported mode")
}

func (a *App) renderDir(inputDirpath string, outputDirpath string, excludedPaths, excludeSingularGlobs []string, data map[string]any) error {
	err := filepath.WalkDir(inputDirpath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(inputDirpath, path)
		if err != nil {
			return fmt.Errorf("get relative path: %s", err)
		}

		var output io.Writer = os.Stdout
		var closer func()
		if slices.Contains(excludedPaths, filepath.Join(inputDirpath, relPath)) {
			return nil
		} else {
			for _, sg := range excludeSingularGlobs {
				if glob.MustCompile(sg).Match(relPath) || glob.MustCompile(sg).Match(filepath.Base(relPath)) {
					return nil
				}
			}
		}
		if len(outputDirpath) > 0 {
			output, closer, err = createOutputFileWithDir(filepath.Join(outputDirpath, relPath))
			if err != nil {
				return err
			}
			defer closer()
		}
		if err := a.renderFile(path, output, data); err != nil {
			return fmt.Errorf("render file %q: %s", path, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("walk directory %q: %s", inputDirpath, err)
	}

	return nil
}

func (a *App) renderFile(inputFilepath string, output io.Writer, data map[string]any) error {
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

func createOutputFileWithDir(outputFilepath string) (io.Writer, func(), error) {
	outputDirpath := filepath.Dir(outputFilepath)
	if err := os.MkdirAll(outputDirpath, os.ModePerm); err != nil {
		return nil, nil, fmt.Errorf("create output directory %s: %s", outputDirpath, err)
	}

	outputFile, err := os.Create(outputFilepath)
	if err != nil {
		return nil, nil, fmt.Errorf("create output file %s: %s", outputFilepath, err)
	}

	return outputFile, func() { outputFile.Close() }, nil
}
