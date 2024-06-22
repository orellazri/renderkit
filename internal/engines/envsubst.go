package engines

import (
	"fmt"
	"io"
	"os"

	"github.com/a8m/envsubst"
)

type EnvsubstEngine struct{}

func (e *EnvsubstEngine) RenderFile(file string, w io.Writer, data map[string]any) error {
	buf, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = envsubstRender(w, buf, data)
	if err != nil {
		return err
	}

	return nil
}

func (e *EnvsubstEngine) Render(r io.Reader, w io.Writer, data map[string]any) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	err = envsubstRender(w, buf, data)
	if err != nil {
		return err
	}

	return nil
}

func envsubstRender(w io.Writer, input []byte, data map[string]any) error {
	// Set environment variables. This is necessary because envsubst does not allow directly passing data as environment variables.
	for k, v := range data {
		os.Setenv(k, fmt.Sprintf("%v", v))
	}

	tpl, err := envsubst.Bytes(input)
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
