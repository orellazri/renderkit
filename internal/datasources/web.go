package datasources

import (
	"bytes"
	"io"
	"mime"
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

	ct := res.Header.Get("Content-Type")
	mt, _, _ := mime.ParseMediaType(ct)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)

	if mt == "application/json" {
		data, err = DecodeJson(r, data)
		if err != nil {
			return nil, err
		}
	} else if mt == "application/toml" {
		data, err = DecodeToml(r, data)
		if err != nil {
			return nil, err
		}
	} else if mt == "application/yaml" || mt == "text/yaml" || mt == "text/x-yaml" || mt == "application/x-yaml" {
		data, err = DecodeYaml(r, data)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}
