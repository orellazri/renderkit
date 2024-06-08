package engine

import "io"

type Engine interface {
	Render(file string, w io.Writer, data map[string]any) error
}

var EnginesMap = map[string]Engine{
	"jet":   &JetEngine{},
	"jinja": &JinjaEngine{},
}
