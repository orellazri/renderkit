package datasources

import (
	"bytes"
	"io"
	"net/http"
)

type WebFileDatasource struct {
	uri string
}

func NewWebFileDatasource(uri string) *WebFileDatasource {
	return &WebFileDatasource{uri}
}

func (ds *WebFileDatasource) Load() (map[string]any, error) {
	data := make(map[string]any)

	res, err := http.Get(ds.uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// b = []byte(`name = "toml"`)
	// b = []byte(`{"name": "json"}`)
	// b = []byte(`name: yaml`)
	r := bytes.NewReader(b)

	if IsJson(b) {
		data, err = DecodeJson(r, data)
		if err != nil {
			return nil, err
		}
	} else if IsToml(b) {
		data, err = DecodeToml(r, data)
		if err != nil {
			return nil, err
		}
	} else if IsYaml(b) {
		data, err = DecodeYaml(r, data)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}
