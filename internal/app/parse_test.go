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
	tmpDir := t.TempDir()
	file, err := os.Create(filepath.Join(tmpDir, "ds.yaml"))
	require.NoError(t, err)
	url, err := url.Parse(file.Name())
	require.NoError(t, err)
	ds, err := a.createDatasourceFromURL(url)
	require.NoError(t, err)
	require.IsType(t, &datasources.YamlDatasource{}, ds)
}

func TestCreateJsonDatasourceFromURL(t *testing.T) {
	a := &App{}
	tmpDir := t.TempDir()
	file, err := os.Create(filepath.Join(tmpDir, "ds.json"))
	require.NoError(t, err)
	url, err := url.Parse(file.Name())
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
	datasources := []string{"/tmp/ds.yaml", "/tmp/ds.json", "/tmp/ds.toml"}
	expectedUrls := []*url.URL{
		{Path: "/tmp/ds.yaml"},
		{Path: "/tmp/ds.json"},
		{Path: "/tmp/ds.toml"},
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
	ds3File, err := os.Create(filepath.Join(tmpDir, "ds3.toml"))
	require.NoError(t, err)

	_, err = ds1File.WriteString("key1: value1")
	require.NoError(t, err)
	_, err = ds2File.WriteString(`{"key2": "value2"}`)
	require.NoError(t, err)
	_, err = ds3File.WriteString(`key3 = "value3"`)
	require.NoError(t, err)
	extraData := []string{"key4=value4"}

	a := &App{}
	datasourceUrls := []*url.URL{
		{Path: ds1File.Name()},
		{Path: ds2File.Name()},
		{Path: ds3File.Name()},
	}
	expectedData := map[string]any{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
	}
	data, err := a.loadDatasources(datasourceUrls, extraData, false)
	require.NoError(t, err)
	require.Equal(t, expectedData, data)
}

func TestCompileGlob(t *testing.T) {
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

func TestCompileInvalidGlob(t *testing.T) {
	app := &App{}
	_, err := app.compileGlob("[a-z")
	require.Error(t, err)
}

func TestAggregateExcludeFiles(t *testing.T) {
	app := &App{}

	tmpDir := t.TempDir()
	for _, file := range []string{"1.txt", "2.txt", "3.txt", "4.txt"} {
		_, err := os.Create(filepath.Join(tmpDir, file))
		require.NoError(t, err)
	}

	excludeFilesGlobs := []string{filepath.Join(tmpDir, "[1-2].txt"), filepath.Join(tmpDir, "3*.txt")}

	aggregatedExcludeFiles, err := app.aggregateExcludeFiles(excludeFilesGlobs)
	require.NoError(t, err)

	require.ElementsMatch(t, []string{
		filepath.Join(tmpDir, "1.txt"),
		filepath.Join(tmpDir, "2.txt"),
		filepath.Join(tmpDir, "3.txt"),
	}, aggregatedExcludeFiles)
}
