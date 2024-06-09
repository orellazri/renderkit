package app

import (
	"net/url"
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
