package engines

import (
	"io"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type GoTemplatesEngine struct{}

func (e *GoTemplatesEngine) RenderFile(file string, w io.Writer, data map[string]any) error {
	tpl, err := template.New(filepath.Base(file)).Funcs(sprig.FuncMap()).ParseFiles(file)
	if err != nil {
		return err
	}

	err = tpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func (e *GoTemplatesEngine) Render(r io.Reader, w io.Writer, data map[string]any) error {
	contents, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	tpl, err := template.New("template").Funcs(sprig.FuncMap()).Parse(string(contents))
	if err != nil {
		return err
	}

	err = tpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
