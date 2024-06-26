package engines

import (
	"fmt"
	"io"
	"os"

	"github.com/a8m/envsubst"
)

type EnvsubstEngine struct{}

func (e *EnvsubstEngine) RenderFile(file string, w io.Writer, data map[string]any) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return e.Render(f, w, data)
}

func (e *EnvsubstEngine) Render(r io.Reader, w io.Writer, data map[string]any) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	// Set environment variables. This is necessary because envsubst does not allow directly passing data as environment variables.
	for k, v := range data {
		os.Setenv(k, fmt.Sprintf("%v", v))
	}

	tpl, err := envsubst.Bytes(buf)
	if err != nil {
		return err
	}

	_, err = w.Write(tpl)
	if err != nil {
		return err
	}

	// Unset environment variables.
	for k := range data {
		os.Unsetenv(k)
	}

	return nil
}
