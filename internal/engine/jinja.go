package engine

import (
	"fmt"
	"io"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

type JinjaEngine struct{}

func (e *JinjaEngine) Render(file string, w io.Writer, data map[string]any) error {
	template, err := gonja.FromFile(file)
	if err != nil {
		return fmt.Errorf("parse file %s: %s", file, err)
	}

	dataCtx := exec.NewContext(data)
	if err := template.Execute(w, dataCtx); err != nil {
		return fmt.Errorf("execute template: %s", err)
	}

	return nil
}
