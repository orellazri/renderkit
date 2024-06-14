package datasources

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-envparse"
)

type EnvDatasource struct {
	variable string
}

func NewEnvDatasource(variable string) *EnvDatasource {
	return &EnvDatasource{variable}
}

func (ds *EnvDatasource) Load() (map[string]any, error) {
	data := make(map[string]any)

	if ds.variable == "" { // If no variable is provided, we use all environment variables
		r := bytes.NewReader([]byte(strings.Join(os.Environ(), "\n")))
		env, err := envparse.Parse(r)
		if err != nil {
			return nil, fmt.Errorf("parse environment variables: %s", err)
		}
		for k, v := range env {
			data[k] = v
		}
	} else {
		value := os.Getenv(ds.variable)
		if len(value) == 0 {
			return nil, fmt.Errorf("environment variable %q not found", ds.variable)
		}
		data[ds.variable] = value
	}

	return data, nil
}
