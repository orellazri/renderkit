package engine

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "test.txt")
	require.NoError(t, err)

	_, err = file.WriteString("Hello, {{ Name }}! You are {{ Age }} years old.")
	require.NoError(t, err)

	engine := &JetEngine{}
	writer := &bytes.Buffer{}
	err = engine.Render(file.Name(), writer, map[string]any{
		"Name": "John",
		"Age":  20,
	})
	require.NoError(t, err)
	require.Equal(t, "Hello, John! You are 20 years old.", writer.String())
}

func TestRenderWithExtends(t *testing.T) {
	dir := t.TempDir()
	baseFile, err := os.CreateTemp(dir, "base.txt")
	require.NoError(t, err)
	_, err = baseFile.WriteString(`
Contents:
{{ block contents() }}{{ end }}
`)
	require.NoError(t, err)

	childFile, err := os.CreateTemp(dir, "child.txt")
	require.NoError(t, err)
	_, err = childFile.WriteString(fmt.Sprintf(`
{{ extends %q }}
{{ block contents() }}
File contents are here
{{ end }}
`, baseFile.Name()))
	require.NoError(t, err)

	engine := &JetEngine{}
	writer := &bytes.Buffer{}
	err = engine.Render(childFile.Name(), writer, nil)
	require.NoError(t, err)
	require.Equal(t, `
Contents:

File contents are here

`, writer.String())
}
