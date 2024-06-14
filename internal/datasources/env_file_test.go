package datasources

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-envparse"
	"github.com/stretchr/testify/require"
)

func TestEnvFileLoadFromFile(t *testing.T) {
	var expectedData = map[string]any{}
	envVars := []string{"key1=value1", "key2=5"}
	r := bytes.NewReader([]byte(strings.Join(envVars, "\n")))
	env, err := envparse.Parse(r)
	require.NoError(t, err)
	for k, v := range env {
		expectedData[k] = v
	}

	dir := t.TempDir()
	file, err := os.CreateTemp(dir, ".env")
	require.NoError(t, err)
	_, err = file.WriteString(`
key1=value1
key2=5`)
	require.NoError(t, err)

	ds := NewEnvFileDatasource(file.Name())
	data, err := ds.Load()
	require.NoError(t, err)
	require.Equal(t, expectedData, data)
}
