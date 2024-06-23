package datasources

import (
	"io"
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
		return nil, err
	}
	defer file.Close()

	data := make(map[string]any)
	data, err = DecodeYaml(file, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func DecodeYaml(r io.Reader, data map[string]any) (map[string]any, error) {
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func IsYaml(content []byte) bool {
	return yaml.Unmarshal(content, &yaml.Node{}) == nil
}
