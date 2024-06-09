package datasource

import "strings"

type EnvDatasource struct {
	envVars []string
}

func NewEnvDatasource(envVars []string) *EnvDatasource {
	return &EnvDatasource{envVars}
}

func (ds *EnvDatasource) Load() map[string]any {
	data := make(map[string]any)
	for _, envVar := range ds.envVars {
		envVar := strings.Split(envVar, "=")
		data[envVar[0]] = envVar[1]
	}

	return data
}
