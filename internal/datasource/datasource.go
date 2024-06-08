package datasource

import (
	"fmt"
	"net/url"
	"path/filepath"
)

type Datasource interface {
	Load() (map[string]any, error)
}

func CreateDatasourceFromURL(url *url.URL) (Datasource, error) {
	switch url.Scheme {
	case "", "file":
		switch filepath.Ext(url.Path) {
		case ".yaml", ".yml":
			return NewYamlDatasource(url.Path), nil
		case ".json":
			return NewJsonDatasource(url.Path), nil
		default:
			return nil, fmt.Errorf("unsupported file extension: %s", filepath.Ext(url.Path))
		}
	default:
		return nil, fmt.Errorf("scheme not supported: %s", url.Scheme)
	}
}
