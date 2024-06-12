package datasources

import (
	"bytes"
	"fmt"
	"io"
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
	var r io.Reader

	if ds.filepath == "" { // If no file is provided, we use the current environment variables
		r = bytes.NewReader([]byte(strings.Join(os.Environ(), "\n")))
	} else {
		f, err := os.Open(ds.filepath)
		if err != nil {
			return nil, fmt.Errorf("open file %s: %s", ds.filepath, err)
		}
		defer f.Close()
		r = f
	}

	env, err := envparse.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parse environment variables: %s", err)
	}
	for k, v := range env {
		data[k] = v
	}

	return data, nil
}
