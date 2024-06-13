package app

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/orellazri/renderkit/internal/datasources"
	"github.com/stretchr/testify/require"
)

func TestCreateYamlDatasourceFromURL(t *testing.T) {
	a := &App{}
	url, err := url.Parse("/tmp/ds.yaml")
	require.NoError(t, err)
	ds, err := a.createDatasourceFromURL(url)
	require.NoError(t, err)
	require.IsType(t, &datasources.YamlDatasource{}, ds)
}

func TestCreateJsonDatasourceFromURL(t *testing.T) {
	a := &App{}
	url, err := url.Parse("/tmp/ds.json")
	require.NoError(t, err)
	ds, err := a.createDatasourceFromURL(url)
	require.NoError(t, err)
	require.IsType(t, &datasources.JsonDatasource{}, ds)
}

func TestCreateInvalidDatasourceFromURL(t *testing.T) {
	a := &App{}

	// Invalid extension
	url, err := url.Parse("/tmp/ds.nothing")
	require.NoError(t, err)
	_, err = a.createDatasourceFromURL(url)
	require.Error(t, err)

	// Invalid scheme
	url, err = url.Parse("nothing:///tmp/ds.yaml")
	require.NoError(t, err)
	_, err = a.createDatasourceFromURL(url)
	require.Error(t, err)
}

func TestParseDatasourceUrls(t *testing.T) {
	a := &App{}
	datasources := []string{"/tmp/ds.yaml", "/tmp/ds.json"}
	expectedUrls := []*url.URL{
		{Path: "/tmp/ds.yaml"},
		{Path: "/tmp/ds.json"},
	}

	urls, err := a.parseDatasourceUrls(datasources)
	require.NoError(t, err)
	require.Equal(t, expectedUrls, urls)
}

func TestLoadDatasources(t *testing.T) {
	tmpDir := t.TempDir()
	ds1File, err := os.Create(filepath.Join(tmpDir, "ds1.yaml"))
	require.NoError(t, err)
	ds2File, err := os.Create(filepath.Join(tmpDir, "ds2.json"))
	require.NoError(t, err)

	_, err = ds1File.WriteString("key1: value1")
	require.NoError(t, err)
	_, err = ds2File.WriteString(`{"key2": "value2"}`)
	require.NoError(t, err)
	extraData := []string{"key3=value3"}

	a := &App{}
	datasourceUrls := []*url.URL{
		{Path: ds1File.Name()},
		{Path: ds2File.Name()},
	}
	expectedData := map[string]any{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	data, err := a.loadDatasources(datasourceUrls, extraData, false)
	require.NoError(t, err)
	require.Equal(t, expectedData, data)
}

func TestCompileInputGlob(t *testing.T) {
	app := &App{}

	tmpDir := t.TempDir()
	_, err := os.Create(filepath.Join(tmpDir, "input.txt"))
	require.NoError(t, err)

	files, err := app.compileGlob(fmt.Sprintf("%s/*.txt", tmpDir))
	require.NoError(t, err)
	require.Equal(t, []string{filepath.Join(tmpDir, "input.txt")}, files)

	files, err = app.compileGlob(fmt.Sprintf("%s/*", tmpDir))
	require.NoError(t, err)
	require.Equal(t, []string{filepath.Join(tmpDir, "input.txt")}, files)

	files, err = app.compileGlob(fmt.Sprintf("%s/**", tmpDir))
	require.NoError(t, err)
	require.Equal(t, []string{filepath.Join(tmpDir, "input.txt")}, files)

	tmpSubdir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(tmpSubdir, 0755)
	require.NoError(t, err)
	_, err = os.Create(filepath.Join(tmpSubdir, "input2.txt"))
	require.NoError(t, err)

	files, err = app.compileGlob(fmt.Sprintf("%s/**", tmpDir))
	require.NoError(t, err)
	require.ElementsMatch(t, []string{
		filepath.Join(tmpDir, "input.txt"),
		filepath.Join(tmpSubdir, "input2.txt"),
	}, files)
}

func TestCompileInputInvalidGlob(t *testing.T) {
	app := &App{}
	_, err := app.compileGlob("[a-z")
	require.Error(t, err)
}

func TestFilteredInputFiles(t *testing.T) {
	app := &App{}

	tmpDir := t.TempDir()
	for _, file := range []string{"1.txt", "2.txt", "3.txt", "4.txt"} {
		_, err := os.Create(filepath.Join(tmpDir, file))
		require.NoError(t, err)
	}

	excludeFilesGlobs := []string{filepath.Join(tmpDir, "[1-2].txt"), filepath.Join(tmpDir, "3*.txt")}

	aggregatedExcludeFiles, err := app.aggregateExcludeFiles(excludeFilesGlobs)
	require.NoError(t, err)

	inputFiles, err := app.compileGlob(fmt.Sprintf("%s/%s", tmpDir, "*.txt"))
	require.NoError(t, err)
	inputFiles = app.excludeFilesFromInput(inputFiles, aggregatedExcludeFiles)

	require.Equal(t, []string{filepath.Join(tmpDir, "4.txt")}, inputFiles)
}

func TestExcludeFilesFromInput(t *testing.T) {
	app := &App{}

	expectedData := []string{"3.txt"}
	inputFiles := []string{"1.txt", "2.txt", "3.txt"}
	excludeFiles := []string{"1.txt", "2.txt"}

	filteredFiles := app.excludeFilesFromInput(inputFiles, excludeFiles)

	require.Equal(t, expectedData, filteredFiles)
}
