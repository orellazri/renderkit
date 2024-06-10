package engine

import (
	"fmt"
	"io"
	"path/filepath"
	"text/template"
)

type GoTemplatesEngine struct{}

func (e *GoTemplatesEngine) Render(file string, w io.Writer, data map[string]any) error {
	tpl, err := template.New(filepath.Base(file)).ParseFiles(file)
	if err != nil {
		return fmt.Errorf("parse %q: %s", file, err)
	}

	err = tpl.Execute(w, data)
	if err != nil {
		return fmt.Errorf("execute %q: %s", file, err)
	}

	return nil
}
