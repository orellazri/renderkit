package datasource

import "io"

type Datasource interface {
	Load(r io.Reader) (map[string]any, error)
}
