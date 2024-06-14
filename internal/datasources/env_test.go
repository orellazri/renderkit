package datasources

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-envparse"
	"github.com/stretchr/testify/require"
)

func TestEnvLoadFromEnvironment(t *testing.T) {
	var expectedData = map[string]any{}
	r := bytes.NewReader([]byte(strings.Join(os.Environ(), "\n")))
	env, err := envparse.Parse(r)
	require.NoError(t, err)
	for k, v := range env {
		expectedData[k] = v
	}

	ds := NewEnvDatasource("")
	data, err := ds.Load()
	require.NoError(t, err)
	require.Equal(t, expectedData, data)
}

func TestEnvDatasourceLoadVariable(t *testing.T) {
	expectedData := map[string]any{
		"RENDERKIT_VAR2": "value2",
	}

	os.Setenv("RENDERKIT_VAR1", "value1")
	os.Setenv("RENDERKIT_VAR2", "value2")

	ds := NewEnvDatasource("RENDERKIT_VAR2")
	data, err := ds.Load()
	require.NoError(t, err)
	require.Equal(t, expectedData, data)
}

func TestEnvDatasourceLoadVariableNotFound(t *testing.T) {
	ds := NewEnvDatasource("RENDERKIT_NON_EXISTING_ENV_VAR")
	data, err := ds.Load()
	require.Error(t, err)
	require.Nil(t, data)
}
