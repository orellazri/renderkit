package datasource

import (
	"io"

	"gopkg.in/yaml.v3"
)

type YamlDatasource struct{}

func NewYamlDatasource() *YamlDatasource {
	return &YamlDatasource{}
}

func (ds *YamlDatasource) Load(r io.Reader) (map[string]any, error) {
	data := make(map[string]any)
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
