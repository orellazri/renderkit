package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Create input files
	inputDir := t.TempDir()
	inputFile1, err := os.Create(filepath.Join(inputDir, "file1.txt"))
	require.NoError(t, err)
	_, err = inputFile1.WriteString(`Hello, my name is {{ .Name }}. I am {{ .Age }} years old.`)
	require.NoError(t, err)
	inputFile2, err := os.Create(filepath.Join(inputDir, "file2.txt"))
	require.NoError(t, err)
	_, err = inputFile2.WriteString(`I like {{ .Hobby }}.`)
	require.NoError(t, err)

	// Create output directory
	outputDir := t.TempDir()

	// Create datasource files
	datasourceDir := t.TempDir()
	datasource1File, err := os.Create(filepath.Join(datasourceDir, "ds.yaml"))
	require.NoError(t, err)
	_, err = datasource1File.WriteString(`
Name: John
Age: 30`)
	require.NoError(t, err)
	datasource2File, err := os.Create(filepath.Join(datasourceDir, "ds2.json"))
	require.NoError(t, err)
	_, err = datasource2File.WriteString(`{"Hobby": "swimming"}`)
	require.NoError(t, err)

	// Run the app
	app := NewApp()
	err = app.Run([]string{
		"",
		"--input-dir", inputDir,
		"--output-dir", outputDir,
		"--datasource", datasource1File.Name(),
		"--datasource", datasource2File.Name(),
		"--engine", "gotemplates",
	})
	require.NoError(t, err)

	// Check the output files
	outputFile1 := filepath.Join(outputDir, "file1.txt")
	outputContent1, err := os.ReadFile(outputFile1)
	require.NoError(t, err)
	require.Equal(t, "Hello, my name is John. I am 30 years old.", string(outputContent1))

	outputFile2 := filepath.Join(outputDir, "file2.txt")
	outputContent2, err := os.ReadFile(outputFile2)
	require.NoError(t, err)
	require.Equal(t, "I like swimming.", string(outputContent2))
}
