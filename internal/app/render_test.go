package app

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/orellazri/renderkit/internal/engine"
	"github.com/stretchr/testify/require"
)

func TestRenderFileToFile(t *testing.T) {
	dir := t.TempDir()
	inputFile := filepath.Join(dir, "input.txt")
	outputFile := filepath.Join(dir, "output.txt")
	err := os.WriteFile(inputFile, []byte("Hello, {{ .Name }}!"), os.ModePerm)
	require.NoError(t, err)

	app := &App{
		engine: &engine.GoTemplatesEngine{},
	}
	err = app.renderFileToFile(inputFile, outputFile, map[string]any{
		"Name": "John",
	})
	require.NoError(t, err)
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	require.Equal(t, "Hello, John!", string(content))
}

func TestRenderFilesToDir(t *testing.T) {
	dir := t.TempDir()
	inputFiles := []string{
		filepath.Join(dir, "input1.txt"),
		filepath.Join(dir, "input2.txt"),
	}
	for _, inputFile := range inputFiles {
		err := os.WriteFile(inputFile, []byte("Hello, {{ .Name }}!"), os.ModePerm)
		require.NoError(t, err)
	}
	outputDir := filepath.Join(dir, "output")

	app := &App{
		engine: &engine.GoTemplatesEngine{},
	}
	err := app.renderFilesToDir(inputFiles, outputDir, map[string]any{
		"Name": "John",
	})
	require.NoError(t, err)
	outputFiles, err := os.ReadDir(outputDir)
	require.NoError(t, err)
	require.Len(t, outputFiles, 2)
	for _, outputFile := range outputFiles {
		content, err := os.ReadFile(filepath.Join(outputDir, outputFile.Name()))
		require.NoError(t, err)
		expectedContent := fmt.Sprintf("Hello, %s!", "John")
		require.Equal(t, expectedContent, string(content))
	}
}

func TestRenderDirToDir(t *testing.T) {
	dir := t.TempDir()
	inputDir := filepath.Join(dir, "input")
	outputDir := filepath.Join(dir, "output")
	err := os.MkdirAll(inputDir, os.ModePerm)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(inputDir, "file1.txt"), []byte("Hello, {{ .Name }}!"), os.ModePerm)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(inputDir, "file2.txt"), []byte("Goodbye, {{ .Name }}!"), os.ModePerm)
	require.NoError(t, err)

	app := &App{
		engine: &engine.GoTemplatesEngine{},
	}

	err = app.renderDirToDir(inputDir, outputDir, map[string]any{
		"Name": "John",
	})
	require.NoError(t, err)
	outputFiles, err := os.ReadDir(outputDir)
	require.NoError(t, err)
	require.Len(t, outputFiles, 2)
	for _, outputFile := range outputFiles {
		content, err := os.ReadFile(filepath.Join(outputDir, outputFile.Name()))
		require.NoError(t, err)

		switch outputFile.Name() {
		case "file1.txt":
			require.Equal(t, "Hello, John!", string(content))
		case "file2.txt":
			require.Equal(t, "Goodbye, John!", string(content))
		}
	}
}
