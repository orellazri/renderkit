package app

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/orellazri/renderkit/internal/datasource"
	"github.com/stretchr/testify/require"
)

func TestCreateYamlDatasourceFromURL(t *testing.T) {
	a := &App{}
	url, err := url.Parse("/tmp/ds.yaml")
	require.NoError(t, err)
	ds, err := a.createDatasourceFromURL(url)
	require.NoError(t, err)
	require.IsType(t, &datasource.YamlDatasource{}, ds)
}

func TestCreateJsonDatasourceFromURL(t *testing.T) {
	a := &App{}
	url, err := url.Parse("/tmp/ds.json")
	require.NoError(t, err)
	ds, err := a.createDatasourceFromURL(url)
	require.NoError(t, err)
	require.IsType(t, &datasource.JsonDatasource{}, ds)
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

func TestSetMode(t *testing.T) {
	app := &App{}

	err := app.setMode([]string{"file1"}, "")
	require.NoError(t, err)
	require.Equal(t, ModeFileToFile, app.mode)

	err = app.setMode([]string{"file1", "file2"}, "")
	require.NoError(t, err)
	require.Equal(t, ModeFilesToDir, app.mode)

	err = app.setMode([]string{}, "dir1")
	require.NoError(t, err)
	require.Equal(t, ModeDirToDir, app.mode)

	err = app.setMode([]string{}, "")
	require.Error(t, err)
}

func TestSetEngine(t *testing.T) {
	app := &App{}

	for engName, engInterface := range enginesMap {
		err := app.setEngine(engName)
		require.NoError(t, err)
		require.Equal(t, engInterface, app.engine)
	}

	err := app.setEngine("unsupportedEngine")
	require.Error(t, err)
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

	a := &App{}
	datasourceUrls := []*url.URL{
		{Path: ds1File.Name()},
		{Path: ds2File.Name()},
	}
	expectedData := map[string]any{
		"key1": "value1",
		"key2": "value2",
	}
	data, err := a.loadDatasources(datasourceUrls, false)
	require.NoError(t, err)
	require.Equal(t, expectedData, data)
}
