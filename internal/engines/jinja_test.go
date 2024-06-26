package engines

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJinjaRenderFile(t *testing.T) {
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "test.txt")
	require.NoError(t, err)

	_, err = file.WriteString("Hello, {{ Name }}! You are {{ Age }} years old.")
	require.NoError(t, err)

	engine := &JinjaEngine{}
	writer := &bytes.Buffer{}
	err = engine.RenderFile(file.Name(), writer, map[string]any{
		"Name": "John",
		"Age":  20,
	})
	require.NoError(t, err)
	require.Equal(t, "Hello, John! You are 20 years old.", writer.String())
}

func TestJinjaRenderFileAdvanced(t *testing.T) {
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "test.txt")
	require.NoError(t, err)

	_, err = file.WriteString(`
{% macro foo(v) -%}
	Version is: {{ v }}
{%- endmacro %}

{{ foo(version) }}`)
	require.NoError(t, err)

	engine := &JinjaEngine{}
	writer := &bytes.Buffer{}
	err = engine.RenderFile(file.Name(), writer, map[string]any{
		"version": "1.2.3",
	})

	require.NoError(t, err)
	require.Equal(t, `


Version is: 1.2.3`, writer.String())
}

func TestJinjaRenderReader(t *testing.T) {
	engine := &JinjaEngine{}
	writer := &bytes.Buffer{}
	err := engine.Render(bytes.NewBufferString("Hello, {{ Name }}! You are {{ Age }} years old."),
		writer,
		map[string]any{
			"Name": "John",
			"Age":  20,
		})
	require.NoError(t, err)
	require.Equal(t, "Hello, John! You are 20 years old.", writer.String())
}
