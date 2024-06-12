package engines

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJinjaRender(t *testing.T) {
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "test.txt")
	require.NoError(t, err)

	_, err = file.WriteString("Hello, {{ Name }}! You are {{ Age }} years old.")
	require.NoError(t, err)

	engine := &JinjaEngine{}
	writer := &bytes.Buffer{}
	err = engine.Render(file.Name(), writer, map[string]any{
		"Name": "John",
		"Age":  20,
	})
	require.NoError(t, err)
	require.Equal(t, "Hello, John! You are 20 years old.", writer.String())
}

func TestJinjaRenderAdvanced(t *testing.T) {
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
	err = engine.Render(file.Name(), writer, map[string]any{
		"version": "1.2.3",
	})

	require.NoError(t, err)
	require.Equal(t, `


Version is: 1.2.3`, writer.String())
}
