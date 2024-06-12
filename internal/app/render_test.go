package app

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/orellazri/renderkit/internal/engines"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
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
		engine: &engines.GoTemplatesEngine{},
	}
	err := app.render(inputFiles, outputDir, map[string]any{
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
