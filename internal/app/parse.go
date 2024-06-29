package app

import (
	"fmt"
	"io"
	"mime"
	"net/http"
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
		ds, f, err := a.createDatasourceFromURL(url)
		if err != nil {
			return nil, fmt.Errorf("create datasource %q: %s", url, err)
		}
		if f != nil {
			defer f.Close()
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

func (a *App) createDatasourceFromURL(url *url.URL) (datasources.Datasource, io.ReadCloser, error) {
	urlWithoutPrefix := strings.TrimPrefix(url.String(), fmt.Sprintf("%s://", url.Scheme))

	switch url.Scheme {
	case "":
		f, err := os.Open(urlWithoutPrefix)
		if err != nil {
			return nil, nil, err
		}
		switch filepath.Ext(urlWithoutPrefix) {
		case ".yaml", ".yml":
			return datasources.NewYamlDatasource(f), f, nil
		case ".json":
			return datasources.NewJsonDatasource(f), f, nil
		case ".toml":
			return datasources.NewTomlDatasource(f), f, nil
		case ".env":
			return datasources.NewEnvFileDatasource(f), f, nil
		default:
			return nil, nil, fmt.Errorf("unsupported file extension: %s", filepath.Ext(urlWithoutPrefix))
		}
	case "env":
		variable := ""
		if url.Host != "" {
			variable = urlWithoutPrefix
		}
		return datasources.NewEnvDatasource(variable), nil, nil
	case "http", "https":
		res, err := http.Get(url.String())
		if err != nil {
			return nil, nil, err
		}

		ct := res.Header.Get("Content-Type")
		mt, _, _ := mime.ParseMediaType(ct)

		var targetDs datasources.Datasource

		switch mt {
		case "application/json":
			targetDs = datasources.NewJsonDatasource(res.Body)
		case "application/toml":
			targetDs = datasources.NewTomlDatasource(res.Body)
		case "application/yaml", "text/yaml", "text/x-yaml", "application/x-yaml":
			targetDs = datasources.NewYamlDatasource(res.Body)
		default:
			return nil, nil, fmt.Errorf("unsupported content type: %s", mt)
		}

		return targetDs, res.Body, nil
	default:
		return nil, nil, fmt.Errorf("scheme not supported: %s", url.Scheme)
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
