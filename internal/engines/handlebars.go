package engines

import (
	"io"

	"github.com/aymerick/raymond"
)

type HandlebarsEngine struct{}

func (e *HandlebarsEngine) RenderFile(file string, w io.Writer, data map[string]any) error {
	tpl, err := raymond.ParseFile(file)
	if err != nil {
		return err
	}

	result, err := tpl.Exec(data)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, result)
	if err != nil {
		return err
	}

	return nil
}

func (e *HandlebarsEngine) Render(r io.Reader, w io.Writer, data map[string]any) error {
	contents, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	tpl, err := raymond.Parse(string(contents))
	if err != nil {
		return err
	}

	result, err := tpl.Exec(data)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, result)
	if err != nil {
		return err
	}

	return nil
}
