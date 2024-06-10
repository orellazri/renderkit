package engine

import (
	"fmt"
	"io"

	"github.com/aymerick/raymond"
)

type HandlebarsEngine struct{}

func (e *HandlebarsEngine) Render(file string, w io.Writer, data map[string]any) error {
	tpl, err := raymond.ParseFile(file)
	if err != nil {
		return fmt.Errorf("parse %q: %s", file, err)
	}

	result, err := tpl.Exec(data)
	if err != nil {
		return fmt.Errorf("render %q: %s", file, err)
	}

	_, err = io.WriteString(w, result)
	if err != nil {
		return fmt.Errorf("write %q: %s", file, err)
	}

	return nil
}
