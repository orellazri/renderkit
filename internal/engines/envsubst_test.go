package engines

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvsubstRenderFile(t *testing.T) {
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "test.txt")
	require.NoError(t, err)

	_, err = file.WriteString("Hello, ${NAME}! You are ${AGE} years old.")
	require.NoError(t, err)

	engine := &EnvsubstEngine{}
	writer := &bytes.Buffer{}
	err = engine.RenderFile(file.Name(), writer, map[string]any{
		"NAME": "John",
		"AGE":  20,
	})
	require.NoError(t, err)
	require.Equal(t, "Hello, John! You are 20 years old.", writer.String())
}

func TestEnvsubstRenderFileAdvanced(t *testing.T) {
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "test.txt")
	require.NoError(t, err)

	_, err = file.WriteString("${GREETING=Hello}, James! Do you know the ${OTHER}?")
	require.NoError(t, err)

	engine := &EnvsubstEngine{}
	writer := &bytes.Buffer{}
	err = engine.RenderFile(file.Name(), writer, map[string]any{
		"OTHER": []string{"a", "b", "c"},
	})

	require.NoError(t, err)
	require.Equal(t, "Hello, James! Do you know the [a b c]?", writer.String())
}

func TestEnvsubstRenderReader(t *testing.T) {
	engine := &EnvsubstEngine{}
	writer := &bytes.Buffer{}
	err := engine.Render(bytes.NewBufferString("Hello, ${NAME}! You are ${AGE} years old."),
		writer,
		map[string]any{
			"NAME": "John",
			"AGE":  20,
		})
	require.NoError(t, err)
	require.Equal(t, "Hello, John! You are 20 years old.", writer.String())
}
