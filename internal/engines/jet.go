package engines

import (
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/CloudyKit/jet/v6"
)

type JetEngine struct{}

func (e *JetEngine) RenderFile(file string, w io.Writer, data map[string]any) error {
	abs, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	renderer := jet.NewSet(jet.NewOSFileSystemLoader("/"))
	tpl, err := renderer.GetTemplate(abs)
	if err != nil {
		return err
	}

	dataMap := jet.VarMap{}
	for key, value := range data {
		dataMap[key] = reflect.ValueOf(value)
	}

	if err := tpl.Execute(w, dataMap, nil); err != nil {
		return err
	}

	return nil
}

func (e *JetEngine) Render(r io.Reader, w io.Writer, data map[string]any) error {
	f, err := os.CreateTemp("", "jet")
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return err
	}

	if err := e.RenderFile(f.Name(), w, data); err != nil {
		return err
	}

	return nil
}
