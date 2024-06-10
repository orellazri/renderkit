package engine

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoTemplatesRender(t *testing.T) {
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "test.txt")
	require.NoError(t, err)

	_, err = file.WriteString("Hello, {{ .Name }}! You are {{ .Age }} years old.")
	require.NoError(t, err)

	engine := &GoTemplatesEngine{}
	writer := &bytes.Buffer{}
	err = engine.Render(file.Name(), writer, map[string]any{
		"Name": "John",
		"Age":  20,
	})
	require.NoError(t, err)
	require.Equal(t, "Hello, John! You are 20 years old.", writer.String())
}

func TestGoTemplatesRenderAdvanced(t *testing.T) {
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "test.txt")
	require.NoError(t, err)

	_, err = file.WriteString(`
{{range $name := .names}}Hi {{$name}}<br>{{end}}`)
	require.NoError(t, err)

	engine := &GoTemplatesEngine{}
	writer := &bytes.Buffer{}
	err = engine.Render(file.Name(), writer, map[string]any{
		"names": []string{"John", "Doe"},
	})

	require.NoError(t, err)
	require.Equal(t, `
Hi John<br>Hi Doe<br>`, writer.String())
}

func TestGoTemplatesRenderWithSprigFunctions(t *testing.T) {
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "test.txt")
	require.NoError(t, err)

	_, err = file.WriteString(`
{{ "hello!" | upper | repeat 5 }}`)
	require.NoError(t, err)

	engine := &GoTemplatesEngine{}
	writer := &bytes.Buffer{}
	err = engine.Render(file.Name(), writer, map[string]any{
		"names": []string{"John", "Doe"},
	})

	require.NoError(t, err)
	require.Equal(t, `
HELLO!HELLO!HELLO!HELLO!HELLO!`, writer.String())
}
