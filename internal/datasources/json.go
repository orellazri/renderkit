package datasources

import (
	"encoding/json"
	"io"
)

type JsonDatasource struct {
	r io.Reader
}

func NewJsonDatasource(r io.Reader) *JsonDatasource {
	return &JsonDatasource{r}
}

func (ds *JsonDatasource) Load() (map[string]any, error) {
	data := make(map[string]any)
	decoder := json.NewDecoder(ds.r)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
