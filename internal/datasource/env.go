package datasource

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
		envVars := os.Environ()
		for _, envVar := range envVars {
			envVar := strings.Split(envVar, "=")
			data[envVar[0]] = envVar[1]
		}
	} else {
		f, err := os.Open(ds.filepath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			envVar := strings.Split(scanner.Text(), "=")
			if len(envVar) != 2 {
				return nil, fmt.Errorf("invalid env var %q in file %s", scanner.Text(), f.Name())
			}
			data[envVar[0]] = envVar[1]
		}
	}
	return data, nil
}
