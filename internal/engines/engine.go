package engines

import "io"

type Engine interface {
	RenderFile(file string, w io.Writer, data map[string]any) error
	Render(r io.Reader, w io.Writer, data map[string]any) error
}
