package datasource

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvLoad(t *testing.T) {
	envVars := []string{
		"key1=value1",
		"KEY2=VALUE2",
	}
	expectedData := map[string]any{
		"key1": "value1",
		"KEY2": "VALUE2",
	}

	ds := NewEnvDatasource(envVars)
	data := ds.Load()
	require.Equal(t, expectedData, data)
}
