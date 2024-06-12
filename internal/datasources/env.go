package datasources

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-envparse"
)

type EnvDatasource struct {
	filepath string
}

func NewEnvDatasource(filepath string) *EnvDatasource {
	return &EnvDatasource{filepath}
}

func (ds *EnvDatasource) Load() (map[string]any, error) {
	data := make(map[string]any)
	if ds.filepath == "" {
		r := bytes.NewReader([]byte(strings.Join(os.Environ(), "\n")))
		env, err := envparse.Parse(r)
		if err != nil {
			return nil, fmt.Errorf("parse environment variables: %s", err)
		}
		for k, v := range env {
			data[k] = v
		}
	} else {
		f, err := os.Open(ds.filepath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		env, err := envparse.Parse(f)
		if err != nil {
			return nil, fmt.Errorf("parse environment variables from file %s: %s", ds.filepath, err)
		}
		for k, v := range env {
			data[k] = v
		}
	}
	return data, nil
}
