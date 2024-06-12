package app

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Create output directory
	outputDir := t.TempDir()

	// Create datasource files
	datasourceDir := t.TempDir()
	datasource1File, err := os.Create(filepath.Join(datasourceDir, "ds.yaml"))
	require.NoError(t, err)
	_, err = datasource1File.WriteString("Name: John")
	require.NoError(t, err)
	datasource2File, err := os.Create(filepath.Join(datasourceDir, "ds2.json"))
	require.NoError(t, err)
	_, err = datasource2File.WriteString(`{"Age": 31.5}`)
	require.NoError(t, err)

	// Define the input syntax for each engine
	inputSyntax := map[string]string{
		"gotemplates": `Hello, my name is {{ .Name }}. I am {{ .Age }} years old.`,
		"jinja":       `Hello, my name is {{ Name }}. I am {{ Age }} years old.`,
		"handlebars":  `Hello, my name is {{ Name }}. I am {{ Age }} years old.`,
		"mustache":    `Hello, my name is {{ Name }}. I am {{ Age }} years old.`,
		"jet":         `Hello, my name is {{ Name }}. I am {{ Age }} years old.`,
	}
	require.Equal(t, len(enginesMap), len(inputSyntax), "all engines must be tested")

	for engine, syntax := range inputSyntax {
		inputDir := t.TempDir()
		inputFile, err := os.Create(filepath.Join(inputDir, "file.txt"))
		require.NoError(t, err)
		_, err = inputFile.WriteString(syntax)
		require.NoError(t, err)

		// Run the app for each engine
		app := NewApp()
		err = app.Run([]string{
			"",
			"--input", fmt.Sprintf("%s/*.txt", inputDir),
			"--output", outputDir,
			"--datasource", datasource1File.Name(),
			"--datasource", datasource2File.Name(),
			"--engine", engine,
		})
		require.NoError(t, err)

		// Check the output files
		outputFile1 := filepath.Join(outputDir, "file.txt")
		outputContent1, err := os.ReadFile(outputFile1)
		require.NoError(t, err)
		require.Equal(t, "Hello, my name is John. I am 31.5 years old.", string(outputContent1))
	}
}
