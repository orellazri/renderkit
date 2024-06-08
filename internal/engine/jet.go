package engine

import (
	"fmt"
	"io"
	"path/filepath"
	"reflect"

	"github.com/CloudyKit/jet/v6"
)

type JetEngine struct{}

func (e *JetEngine) Render(file string, w io.Writer, data map[string]any) error {
	abs, err := filepath.Abs(file)
	if err != nil {
		return fmt.Errorf("get absolute path: %s", err)
	}

	renderer := jet.NewSet(jet.NewOSFileSystemLoader("/"))
	tpl, err := renderer.GetTemplate(abs)
	if err != nil {
		return fmt.Errorf("get template: %s", err)
	}

	dataMap := jet.VarMap{}
	for key, value := range data {
		dataMap[key] = reflect.ValueOf(value)
	}

	if err := tpl.Execute(w, dataMap, nil); err != nil {
		return fmt.Errorf("execute template: %s", err)
	}

	return nil
}
