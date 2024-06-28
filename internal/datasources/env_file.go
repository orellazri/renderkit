package datasources

import (
	"fmt"
	"io"

	"github.com/hashicorp/go-envparse"
)

type EnvFileDatasource struct {
	r io.Reader
}

func NewEnvFileDatasource(r io.Reader) *EnvFileDatasource {
	return &EnvFileDatasource{r}
}

func (ds *EnvFileDatasource) Load() (map[string]any, error) {
	data := make(map[string]any)

	env, err := envparse.Parse(ds.r)
	if err != nil {
		return nil, fmt.Errorf("parse environment variables: %s", err)
	}

	for k, v := range env {
		data[k] = v
	}

	return data, nil
}
