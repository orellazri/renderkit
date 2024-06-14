package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAllEngines(t *testing.T) {
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
			"--input-dir", inputDir,
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

func TestIntegrationInputOutputSubdirsMirrored(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Create input files
	inputDir := t.TempDir()
	inputSubdir1 := filepath.Join(inputDir, "subdir1")
	err := os.Mkdir(inputSubdir1, os.ModePerm)
	require.NoError(t, err)
	inputSubdir2 := filepath.Join(inputDir, "subdir2")
	err = os.Mkdir(inputSubdir2, os.ModePerm)
	require.NoError(t, err)

	inputFiles := []string{
		filepath.Join(inputSubdir1, "file1.txt"),
		filepath.Join(inputSubdir2, "file2.txt"),
	}

	err = os.WriteFile(inputFiles[0], []byte("My name is {{ .Name }}"), os.ModePerm)
	require.NoError(t, err)
	err = os.WriteFile(inputFiles[1], []byte("I am {{ .Age }} years old"), os.ModePerm)
	require.NoError(t, err)

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

	app := NewApp()
	err = app.Run([]string{
		"",
		"--input-dir", inputDir,
		"--output", outputDir,
		"--datasource", datasource1File.Name(),
		"--datasource", datasource2File.Name(),
	})
	require.NoError(t, err)

	// Check the output files
	outputFile1 := filepath.Join(outputDir, "subdir1", "file1.txt")
	outputContent1, err := os.ReadFile(outputFile1)
	require.NoError(t, err)
	require.Equal(t, "My name is John", string(outputContent1))

	outputFile2 := filepath.Join(outputDir, "subdir2", "file2.txt")
	outputContent2, err := os.ReadFile(outputFile2)
	require.NoError(t, err)
	require.Equal(t, "I am 31.5 years old", string(outputContent2))
}
