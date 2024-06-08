package datasource

type Datasource interface {
	Load() (map[string]any, error)
}
