package engines

import (
	"io"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

type JinjaEngine struct{}

func (e *JinjaEngine) RenderFile(file string, w io.Writer, data map[string]any) error {
	tpl, err := gonja.FromFile(file)
	if err != nil {
		return err
	}

	dataCtx := exec.NewContext(data)
	if err := tpl.Execute(w, dataCtx); err != nil {
		return err
	}

	return nil
}

func (e *JinjaEngine) Render(r io.Reader, w io.Writer, data map[string]any) error {
	contents, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	tpl, err := gonja.FromBytes(contents)
	if err != nil {
		return err
	}

	dataCtx := exec.NewContext(data)
	if err := tpl.Execute(w, dataCtx); err != nil {
		return err
	}

	return nil
}
