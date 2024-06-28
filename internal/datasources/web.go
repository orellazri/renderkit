package datasources

import (
	"fmt"
	"mime"
	"net/http"
)

type WebFileDatasource struct {
	url string
}

func NewWebFileDatasource(url string) *WebFileDatasource {
	return &WebFileDatasource{url}
}

func (ds *WebFileDatasource) Load() (map[string]any, error) {
	res, err := http.Get(ds.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	ct := res.Header.Get("Content-Type")
	mt, _, _ := mime.ParseMediaType(ct)

	var targetDs Datasource

	switch mt {
	case "application/json":
		targetDs = NewJsonDatasource(res.Body)
	case "application/toml":
		targetDs = NewTomlDatasource(res.Body)
	case "application/yaml", "text/yaml", "text/x-yaml", "application/x-yaml":
		targetDs = NewYamlDatasource(res.Body)
	default:
		return nil, fmt.Errorf("unsupported content type: %s", mt)
	}

	data, err := targetDs.Load()
	if err != nil {
		return nil, err
	}

	return data, nil
}
