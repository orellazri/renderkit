package datasources

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type TomlDatasource struct {
	filepath string
}

func NewTomlDatasource(filepath string) *TomlDatasource {
	return &TomlDatasource{filepath}
}

func (ds *TomlDatasource) Load() (map[string]any, error) {
	file, err := os.Open(ds.filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make(map[string]any)
	decoder := toml.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}