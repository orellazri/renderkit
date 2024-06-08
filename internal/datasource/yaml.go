package datasource

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type YamlDatasource struct {
	filepath string
}

func NewYamlDatasource(filepath string) *YamlDatasource {
	return &YamlDatasource{filepath}
}

func (ds *YamlDatasource) Load() (map[string]any, error) {
	file, err := os.Open(ds.filepath)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %s", ds.filepath, err)
	}
	defer file.Close()

	data := make(map[string]any)
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
