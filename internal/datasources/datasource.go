package datasources

type Datasource interface {
	Load() (map[string]any, error)
}
