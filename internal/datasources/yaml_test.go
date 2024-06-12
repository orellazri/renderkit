package datasources

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYamlLoad(t *testing.T) {
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "ds.yaml")
	require.NoError(t, err)

	_, err = file.WriteString(`
key1: value1
key2: 5`)
	require.NoError(t, err)
	ds := NewYamlDatasource(file.Name())

	expectedData := map[string]any{
		"key1": "value1",
		"key2": 5,
	}
	data, err := ds.Load()
	require.NoError(t, err)
	require.Equal(t, expectedData, data)
}
