package datasources

import (
	"io"

	"github.com/pelletier/go-toml/v2"
)

type TomlDatasource struct {
	r io.Reader
}

func NewTomlDatasource(r io.Reader) *TomlDatasource {
	return &TomlDatasource{r}
}

func (ds *TomlDatasource) Load() (map[string]any, error) {
	data := make(map[string]any)
	decoder := toml.NewDecoder(ds.r)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
