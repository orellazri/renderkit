package datasources

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvFileLoadFromFile(t *testing.T) {
	envFileData := `
key1=value1
key2=5`
	expectedData := map[string]any{
		"key1": "value1",
		"key2": "5",
	}
	r := strings.NewReader(envFileData)
	ds := NewEnvFileDatasource(r)

	data, err := ds.Load()
	require.NoError(t, err)
	require.Equal(t, expectedData, data)
}
