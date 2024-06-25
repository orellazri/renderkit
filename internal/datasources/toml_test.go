package datasources

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTomlLoad(t *testing.T) {
	tomlData := `
key1 = "value1"
key2 = 5`
	expectedData := map[string]any{
		"key1": "value1",
		"key2": int64(5),
	}
	r := strings.NewReader(tomlData)
	ds := NewTomlDatasource(r)

	data, err := ds.Load()
	require.NoError(t, err)
	require.Equal(t, expectedData, data)
}
