package datasources

import (
	"io"

	"gopkg.in/yaml.v3"
)

type YamlDatasource struct {
	r io.Reader
}

func NewYamlDatasource(r io.Reader) *YamlDatasource {
	return &YamlDatasource{r}
}

func (ds *YamlDatasource) Load() (map[string]any, error) {
	data := make(map[string]any)
	decoder := yaml.NewDecoder(ds.r)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
