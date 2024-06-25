package datasources

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJsonLoad(t *testing.T) {
	jsonData := `
{
	"key1": "value1",
	"key2": 5
}`
	expectedData := map[string]any{
		"key1": "value1",
		"key2": float64(5),
	}
	r := strings.NewReader(jsonData)
	ds := NewJsonDatasource(r)

	data, err := ds.Load()
	require.NoError(t, err)
	require.Equal(t, expectedData, data)
}
