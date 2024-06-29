package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

	// Define the input syntax for each engine
	inputSyntax := map[string]string{
		"envsubst":    "Hello, my name is ${Name}. I am ${Age} years old.",
		"gotemplates": `Hello, my name is {{ .Name }}. I am {{ .Age }} years old.`,
		"handlebars":  `Hello, my name is {{ Name }}. I am {{ Age }} years old.`,
		"jet":         `Hello, my name is {{ Name }}. I am {{ Age }} years old.`,
		"jinja":       `Hello, my name is {{ Name }}. I am {{ Age }} years old.`,
		"mustache":    `Hello, my name is {{ Name }}. I am {{ Age }} years old.`,
	}
	require.Equal(t, len(enginesMap), len(inputSyntax), "all engines must be tested")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err = fmt.Fprint(w, `{"Age": 31.5}`)
		require.NoError(t, err)
	}))
	defer ts.Close()

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
			"--datasource", ts.URL,
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
	inputSubdir3 := filepath.Join(inputDir, "subdir3")
	err = os.Mkdir(inputSubdir3, os.ModePerm)
	require.NoError(t, err)

	inputFiles := []string{
		filepath.Join(inputSubdir1, "file1.txt"),
		filepath.Join(inputSubdir2, "file2.txt"),
		filepath.Join(inputSubdir3, "file3.txt"),
	}

	err = os.WriteFile(inputFiles[0], []byte("My name is {{ .Name }} and I am {{ .Age }} years old"), os.ModePerm)
	require.NoError(t, err)
	err = os.WriteFile(inputFiles[1], []byte("I am {{ .Age }} years old. This file will be excluded."), os.ModePerm)
	require.NoError(t, err)
	err = os.WriteFile(inputFiles[2], []byte("I am {{ .Age }} years old. This file will also be excluded."), os.ModePerm)
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
		"--exclude", fmt.Sprintf("%s/*2.txt", inputSubdir2),
		"--exclude", "*3.txt",
		"--output", outputDir,
		"--datasource", datasource1File.Name(),
		"--datasource", datasource2File.Name(),
	})
	require.NoError(t, err)

	// Check the output files
	outputFile1 := filepath.Join(outputDir, "subdir1", "file1.txt")
	outputContent1, err := os.ReadFile(outputFile1)
	require.NoError(t, err)
	require.Equal(t, "My name is John and I am 31.5 years old", string(outputContent1))

	// Check that the excluded files aren't present
	outputFile2 := filepath.Join(outputDir, "subdir2", "file2.txt")
	_, err = os.Stat(outputFile2)
	require.ErrorIs(t, err, os.ErrNotExist)
	outputFile3 := filepath.Join(outputDir, "subdir3", "file3.txt")
	_, err = os.Stat(outputFile3)
	require.ErrorIs(t, err, os.ErrNotExist)
}
