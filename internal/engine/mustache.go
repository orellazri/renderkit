package engine

import (
	"fmt"
	"io"

	"github.com/cbroglie/mustache"
)

type MustacheEngine struct{}

func (e *MustacheEngine) Render(file string, w io.Writer, data map[string]any) error {
	out, err := mustache.RenderFile(file, data)
	if err != nil {
		return fmt.Errorf("render %q: %s", file, err)
	}

	_, err = io.WriteString(w, out)
	if err != nil {
		return fmt.Errorf("write %q: %s", file, err)
	}

	return nil
}
