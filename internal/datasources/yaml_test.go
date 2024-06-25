package datasources

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYamlLoad(t *testing.T) {
	yamlData := `
key1: value1
key2: 5`
	expectedData := map[string]any{
		"key1": "value1",
		"key2": 5,
	}
	r := strings.NewReader(yamlData)
	ds := NewYamlDatasource(r)

	data, err := ds.Load()
	require.NoError(t, err)
	require.Equal(t, expectedData, data)
}
