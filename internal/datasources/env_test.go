package datasources

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvLoad(t *testing.T) {
	expectedData := map[string]any{}
	envVars := os.Environ()
	for _, envVar := range envVars {
		envVar := strings.Split(envVar, "=")
		expectedData[envVar[0]] = envVar[1]
	}

	ds := NewEnvDatasource("")
	data, err := ds.Load()
	require.NoError(t, err)
	require.Equal(t, expectedData, data)

	dir := t.TempDir()
	file, err := os.CreateTemp(dir, ".env")
	require.NoError(t, err)
	_, err = file.WriteString(`key1=value1
key2=5`)
	require.NoError(t, err)

	expectedData = map[string]any{
		"key1": "value1",
		"key2": "5",
	}
	ds = NewEnvDatasource(file.Name())
	data, err = ds.Load()
	require.NoError(t, err)
	require.Equal(t, expectedData, data)
}
