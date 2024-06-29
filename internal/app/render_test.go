package app

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/orellazri/renderkit/internal/engines"
	"github.com/stretchr/testify/require"
)

func TestRenderDir(t *testing.T) {
	dir := t.TempDir()
	inputDir := filepath.Join(dir, "input")
	err := os.Mkdir(inputDir, os.ModePerm)

	require.NoError(t, err)
	inputFiles := []string{
		filepath.Join(inputDir, "input1.txt"),
		filepath.Join(inputDir, "input2.txt"),
	}
	for _, inputFile := range inputFiles {
		err := os.WriteFile(inputFile, []byte("Hello, {{ .Name }}!"), os.ModePerm)
		require.NoError(t, err)
	}
	outputDir := filepath.Join(dir, "output")

	app := &App{
		engine: &engines.GoTemplatesEngine{},
	}
	err = app.render(
		"",
		inputDir,
		"",
		outputDir,
		nil,
		nil,
		map[string]any{
			"Name": "John",
		},
	)
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
func TestRenderDirWithSubpaths(t *testing.T) {
	dir := t.TempDir()
	inputDir := filepath.Join(dir, "input")
	err := os.Mkdir(inputDir, os.ModePerm)
	require.NoError(t, err)

	inputSubdir1 := filepath.Join(inputDir, "subdir1")
	err = os.Mkdir(inputSubdir1, os.ModePerm)
	require.NoError(t, err)
	inputSubdir2 := filepath.Join(inputDir, "subdir2")
	err = os.Mkdir(inputSubdir2, os.ModePerm)
	require.NoError(t, err)

	inputFiles := []string{
		filepath.Join(inputSubdir1, "file1.txt"),
		filepath.Join(inputSubdir2, "file2.txt"),
	}
	for _, inputFile := range inputFiles {
		err := os.WriteFile(inputFile, []byte("Hello!"), os.ModePerm)
		require.NoError(t, err)
	}

	outputDir := filepath.Join(dir, "output")
	app := &App{
		engine: &engines.GoTemplatesEngine{},
	}

	err = app.renderDir(inputDir, outputDir, nil, nil, nil)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(outputDir, "subdir1", "file1.txt"))
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(outputDir, "subdir2", "file2.txt"))
	require.NoError(t, err)
}

func TestRenderFile(t *testing.T) {
	dir := t.TempDir()
	inputFile := filepath.Join(dir, "input.txt")
	err := os.WriteFile(inputFile, []byte("Hello, {{ .Name }}!"), os.ModePerm)
	require.NoError(t, err)
	outputDir := filepath.Join(dir, "output")
	app := &App{
		engine: &engines.GoTemplatesEngine{},
	}
	err = app.render(
		"",
		"",
		inputFile,
		outputDir,
		nil,
		nil,
		map[string]any{
			"Name": "John",
		},
	)
	require.NoError(t, err)
	outputFile := filepath.Join(outputDir, filepath.Base(inputFile))
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	expectedContent := fmt.Sprintf("Hello, %s!", "John")
	require.Equal(t, expectedContent, string(content))
}

func TestRenderFromString(t *testing.T) {
	app := &App{
		engine: &engines.GoTemplatesEngine{},
	}
	input := "Hello, {{ .Name }}!"
	buf := &bytes.Buffer{}
	err := app.renderString(input, buf, map[string]any{
		"Name": "John",
	})
	require.NoError(t, err)
	require.Equal(t, "Hello, John!", buf.String())
}

func TestRenderFromStringToFile(t *testing.T) {
	tmpDir := t.TempDir()
	app := &App{
		engine: &engines.GoTemplatesEngine{},
	}
	err := app.render(
		"Hello, {{ .Name }}!",
		"",
		"",
		tmpDir,
		nil,
		nil,
		map[string]any{
			"Name": "John",
		},
	)
	require.NoError(t, err)
	outputFile := filepath.Join(tmpDir, "renderkit_output")
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	require.Equal(t, "Hello, John!", string(content))
}
