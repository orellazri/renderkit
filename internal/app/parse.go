package app

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/orellazri/renderkit/internal/datasource"
	"github.com/orellazri/renderkit/internal/engine"
)

var enginesMap = map[string]engine.Engine{
	"gotemplates": &engine.GoTemplatesEngine{},
	"jinja":       &engine.JinjaEngine{},
	"handlebars":  &engine.HandlebarsEngine{},
	"mustache":    &engine.MustacheEngine{},
	"jet":         &engine.JetEngine{},
}

func (a *App) parseDatasourceUrls(datasources []string) ([]*url.URL, error) {
	datasourceUrls := make([]*url.URL, len(datasources))
	for i, ds := range datasources {
		url, err := url.Parse(ds)
		if err != nil {
			return nil, fmt.Errorf("invalid url %s: %s", ds, err)
		}
		datasourceUrls[i] = url
	}

	return datasourceUrls, nil
}

func (a *App) loadDatasources(datasourceUrls []*url.URL, allowDuplicateKeys bool) (map[string]any, error) {
	duplicateKeys := []string{} // We keep track of duplicate keys to return a more informative error message
	data := make(map[string]any)
	// Load datasources
	for _, url := range datasourceUrls {
		ds, err := a.createDatasourceFromURL(url)
		if err != nil {
			return nil, fmt.Errorf("create datasource %q: %s", url, err)
		}

		dsData, err := ds.Load()
		if err != nil {
			return nil, fmt.Errorf("load datasource %q: %s", url, err)
		}

		// Merge with data dictionary
		for k, v := range dsData {
			if _, ok := data[k]; ok && !allowDuplicateKeys {
				duplicateKeys = append(duplicateKeys, k)
			}
			data[k] = v
		}
	}

	if len(duplicateKeys) > 0 {
		return nil, fmt.Errorf("duplicate keys found in datasources: %s", strings.Join(duplicateKeys, ", "))
	}

	return data, nil
}

func (a *App) createDatasourceFromURL(url *url.URL) (datasource.Datasource, error) {
	switch url.Scheme {
	case "":
		switch filepath.Ext(url.Path) {
		case ".yaml", ".yml":
			return datasource.NewYamlDatasource(url.Path), nil
		case ".json":
			return datasource.NewJsonDatasource(url.Path), nil
		default:
			return nil, fmt.Errorf("unsupported file extension: %s", filepath.Ext(url.Path))
		}
	case "env":
		if url.Path != "" {
			path := strings.TrimPrefix(url.String(), fmt.Sprintf("%s://", url.Scheme))
			return datasource.NewEnvDatasource(path), nil
		}
		return datasource.NewEnvDatasource(url.Path), nil
	default:
		return nil, fmt.Errorf("scheme not supported: %s", url.Scheme)
	}
}

func (a *App) setMode(input []string, inputDir string) error {
	if len(input) == 1 {
		a.mode = ModeFileToFile
	} else if len(input) > 1 {
		a.mode = ModeFilesToDir
	} else if len(inputDir) > 0 {
		a.mode = ModeDirToDir
	} else {
		return errors.New("unsupported mode")
	}

	return nil
}

func (a *App) setEngine(engineStr string) error {
	e, ok := enginesMap[engineStr]
	if !ok {
		return fmt.Errorf("unsupported engine: %s", engineStr)
	}
	a.engine = e

	return nil
}
