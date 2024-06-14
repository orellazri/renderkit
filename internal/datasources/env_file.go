package datasources

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-envparse"
)

type EnvFileDatasource struct {
	filepath string
}

func NewEnvFileDatasource(filepath string) *EnvFileDatasource {
	return &EnvFileDatasource{filepath}
}

func (ds *EnvFileDatasource) Load() (map[string]any, error) {
	data := make(map[string]any)

	f, err := os.Open(ds.filepath)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %s", ds.filepath, err)
	}
	defer f.Close()

	env, err := envparse.Parse(f)
	if err != nil {
		return nil, fmt.Errorf("parse environment variables from file %s: %s", ds.filepath, err)
	}

	for k, v := range env {
		data[k] = v
	}

	return data, nil
}
