package app

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/goreleaser/fileglob"
	"github.com/orellazri/renderkit/internal/datasources"
	"github.com/orellazri/renderkit/internal/engines"
)

var enginesMap = map[string]engines.Engine{
	"envsubst":    &engines.EnvsubstEngine{},
	"gotemplates": &engines.GoTemplatesEngine{},
	"handlebars":  &engines.HandlebarsEngine{},
	"jet":         &engines.JetEngine{},
	"jinja":       &engines.JinjaEngine{},
	"mustache":    &engines.MustacheEngine{},
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
	urlWithoutPrefix := strings.TrimPrefix(url.String(), fmt.Sprintf("%s://", url.Scheme))

	switch url.Scheme {
	case "":
		switch filepath.Ext(urlWithoutPrefix) {
		case ".yaml", ".yml":
			f, err := os.Open(urlWithoutPrefix)
			if err != nil {
				return nil, err
			}
			// defer f.Close()
			return datasources.NewYamlDatasource(f), nil
		case ".json":
			f, err := os.Open(urlWithoutPrefix)
			if err != nil {
				return nil, err
			}
			// defer f.Close()
			return datasources.NewJsonDatasource(f), nil
		case ".toml":
			f, err := os.Open(urlWithoutPrefix)
			if err != nil {
				return nil, err
			}
			// defer f.Close()
			return datasources.NewTomlDatasource(f), nil
		case ".env":
			f, err := os.Open(urlWithoutPrefix)
			if err != nil {
				return nil, err
			}
			// defer f.Close()
			return datasources.NewEnvFileDatasource(f), nil
		default:
			return nil, fmt.Errorf("unsupported file extension: %s", filepath.Ext(urlWithoutPrefix))
		}
	case "env":
		variable := ""
		if url.Host != "" {
			variable = urlWithoutPrefix
		}
		return datasources.NewEnvDatasource(variable), nil
	case "http", "https":
		return datasources.NewWebFileDatasource(url.String()), nil
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

func (a *App) aggregateExcludeFiles(excludeFiles []string) ([]string, error) {
	var aggregatedExcludeFiles []string
	for _, excludeGlob := range excludeFiles {
		excludeFiles, err := a.compileGlob(excludeGlob)
		if err != nil {
			return nil, fmt.Errorf("compile exclude glob %q: %s", excludeGlob, err)
		}
		aggregatedExcludeFiles = slices.Concat(aggregatedExcludeFiles, excludeFiles)
	}
	slices.Sort(aggregatedExcludeFiles)
	aggregatedExcludeFiles = slices.Compact(aggregatedExcludeFiles)
	return aggregatedExcludeFiles, nil
}
