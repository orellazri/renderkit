package engines

import (
	"io"

	"github.com/cbroglie/mustache"
)

type MustacheEngine struct{}

func (e *MustacheEngine) RenderFile(file string, w io.Writer, data map[string]any) error {
	result, err := mustache.RenderFile(file, data)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, result)
	if err != nil {
		return err
	}

	return nil
}

func (e *MustacheEngine) Render(r io.Reader, w io.Writer, data map[string]any) error {
	contents, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	result, err := mustache.Render(string(contents), data)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, result)
	if err != nil {
		return err
	}

	return nil
}
