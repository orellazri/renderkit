package engines

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlebarsRenderFile(t *testing.T) {
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "test.txt")
	require.NoError(t, err)

	_, err = file.WriteString("Hello, {{ Name }}! You are {{ Age }} years old.")
	require.NoError(t, err)

	engine := &HandlebarsEngine{}
	writer := &bytes.Buffer{}
	err = engine.RenderFile(file.Name(), writer, map[string]any{
		"Name": "John",
		"Age":  20,
	})
	require.NoError(t, err)
	require.Equal(t, "Hello, John! You are 20 years old.", writer.String())
}

func TestHandlebarsRenderFileAdvanced(t *testing.T) {
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "test.txt")
	require.NoError(t, err)

	_, err = file.WriteString(`
{{#names}}Hi {{.}}<br>{{/names}}`)
	require.NoError(t, err)

	engine := &HandlebarsEngine{}
	writer := &bytes.Buffer{}
	err = engine.RenderFile(file.Name(), writer, map[string]any{
		"names": []string{"John", "Doe"},
	})

	require.NoError(t, err)
	require.Equal(t, `
Hi John<br>Hi Doe<br>`, writer.String())
}

func TestHandlebarsRenderReader(t *testing.T) {
	engine := &HandlebarsEngine{}
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
