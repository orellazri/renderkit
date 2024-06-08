package app

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/orellazri/renderkit/internal/datasource"
	"github.com/urfave/cli/v2"
)

func parseDatasourceUrls(cCtx *cli.Context) ([]*url.URL, error) {
	datasourceUrls := make([]*url.URL, len(cCtx.StringSlice("datasource")))
	for i, ds := range cCtx.StringSlice("datasource") {
		url, err := url.Parse(ds)
		if err != nil {
			return nil, fmt.Errorf("invalid url %s: %s", ds, err)
		}
		datasourceUrls[i] = url
	}

	return datasourceUrls, nil
}

func loadDatasources(datasourceUrls []*url.URL, allowDuplicateKeys bool) (map[string]any, error) {
	duplicateKeys := []string{} // We keep track of duplicate keys to return a more informative error message
	data := make(map[string]any)
	for _, url := range datasourceUrls {
		ds, err := datasource.CreateDatasourceFromURL(url)
		if err != nil {
			return nil, fmt.Errorf("create datasource %s: %s", url, err)
		}

		dsData, err := ds.Load()
		if err != nil {
			return nil, fmt.Errorf("load datasource %s: %s", url, err)
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
