package datasources

import (
	"encoding/json"
	"io"
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
		return nil, err
	}
	defer file.Close()

	data := make(map[string]any)
	data, err = DecodeJson(file, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func DecodeJson(r io.Reader, data map[string]any) (map[string]any, error) {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
