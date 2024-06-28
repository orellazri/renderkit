package datasources

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWebFileLoad(t *testing.T) {
	var err error
	dsFiles := []string{"ds.json", "ds.toml", "ds.yaml"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/" + dsFiles[0]:
			w.Header().Set("Content-Type", "application/json")
			_, err = fmt.Fprint(w, `{"key1": "value1", "key2": "value2"}`)
			require.NoError(t, err)
		case "/" + dsFiles[1]:
			w.Header().Set("Content-Type", "application/toml")
			_, err = fmt.Fprint(w, "key1 = \"value1\"\n key2 = \"value2\"")
			require.NoError(t, err)
		case "/" + dsFiles[2]:
			w.Header().Set("Content-Type", "application/yaml")
			_, err = fmt.Fprint(w, "key1: value1\nkey2: value2")
			require.NoError(t, err)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	for _, fileType := range dsFiles {
		ds := NewWebFileDatasource(ts.URL + "/" + fileType)
		data, err := ds.Load()
		require.NoError(t, err)

		expectedData := map[string]any{
			"key1": "value1",
			"key2": "value2",
		}

		require.Equal(t, data, expectedData)
	}
}
