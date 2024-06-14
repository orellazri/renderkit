package app

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/goreleaser/fileglob"
	"github.com/orellazri/renderkit/internal/datasources"
	"github.com/orellazri/renderkit/internal/engines"
)

var enginesMap = map[string]engines.Engine{
	"gotemplates": &engines.GoTemplatesEngine{},
	"jinja":       &engines.JinjaEngine{},
	"handlebars":  &engines.HandlebarsEngine{},
	"mustache":    &engines.MustacheEngine{},
	"jet":         &engines.JetEngine{},
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

func (a *App) loadDatasources(datasourceUrls []*url.URL, extraData []string, allowDuplicateKeys bool) (map[string]any, error) {
	duplicateKeys := []string{} // We keep track of duplicate keys to return a more informative error message
	data := make(map[string]any)

	// Load extra data
	for _, d := range extraData {
		kv := strings.SplitN(d, "=", 2)
		if _, ok := data[kv[0]]; ok && !allowDuplicateKeys {
			duplicateKeys = append(duplicateKeys, kv[0])
		}
		data[kv[0]] = kv[1]
	}

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

func (a *App) createDatasourceFromURL(url *url.URL) (datasources.Datasource, error) {
	switch url.Scheme {
	case "":
		switch filepath.Ext(url.Path) {
		case ".yaml", ".yml":
			return datasources.NewYamlDatasource(url.Path), nil
		case ".json":
			return datasources.NewJsonDatasource(url.Path), nil
		default:
			return nil, fmt.Errorf("unsupported file extension: %s", filepath.Ext(url.Path))
		}
	case "env":
		path := ""
		if url.Host != "" {
			path = strings.TrimPrefix(url.String(), fmt.Sprintf("%s://", url.Scheme))
		}
		return datasources.NewEnvDatasource(path), nil
	default:
		return nil, fmt.Errorf("scheme not supported: %s", url.Scheme)
	}
}

func (a *App) compileGlob(pattern string) ([]string, error) {
	if err := fileglob.ValidPattern(pattern); err != nil {
		return nil, fmt.Errorf("invalid glob pattern: %q", err)
	}
	matches, err := fileglob.Glob(pattern, fileglob.MaybeRootFS)
	if err != nil {
		return nil, fmt.Errorf("glob %q: %s", pattern, err)
	}

	return matches, nil
}

func (a *App) excludeFilesFromInput(inputFiles []string, excludeFiles []string) []string {
	var excludeMap = make(map[string]bool)
	var filtered []string

	for _, v := range excludeFiles {
		excludeMap[v] = true
	}
	for _, v := range inputFiles {
		if _, ok := excludeMap[v]; !ok {
			filtered = append(filtered, v)
		}
	}
	return filtered
}
