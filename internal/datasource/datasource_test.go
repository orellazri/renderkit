package datasource

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateYamlDatasourceFromURL(t *testing.T) {
	url, err := url.Parse("/tmp/ds.yaml")
	require.NoError(t, err)
	ds, err := CreateDatasourceFromURL(url)
	require.NoError(t, err)
	require.IsType(t, &YamlDatasource{}, ds)
	require.Equal(t, "/tmp/ds.yaml", ds.(*YamlDatasource).filepath)
}

func TestCreateJsonDatasourceFromURL(t *testing.T) {
	url, err := url.Parse("/tmp/ds.json")
	require.NoError(t, err)
	ds, err := CreateDatasourceFromURL(url)
	require.NoError(t, err)
	require.IsType(t, &JsonDatasource{}, ds)
	require.Equal(t, "/tmp/ds.json", ds.(*JsonDatasource).filepath)
}
