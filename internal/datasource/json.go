package datasource

import (
	"encoding/json"
	"fmt"
	"os"
)

type JsonDatasource struct {
	filepath string
}

func NewJsonDatasource(filepath string) *JsonDatasource {
	return &JsonDatasource{filepath}
}

func (ds *JsonDatasource) Load() (map[string]any, error) {
	file, err := os.Open(ds.filepath)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %s", ds.filepath, err)
	}
	defer file.Close()

	data := make(map[string]any)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
