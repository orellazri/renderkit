package datasource

import (
	"fmt"
	"net/url"
)

type Datasource interface {
	Load() (map[string]any, error)
}

func CreateDatasourceFromURL(url *url.URL) (Datasource, error) {
	switch url.Scheme {
	case "":
		return NewYamlDatasource(url.Path), nil
	default:
		return nil, fmt.Errorf("scheme not supported: %s", url.Scheme)
	}
}
