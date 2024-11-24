package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateFlagsNoErrors(t *testing.T) {
	app := NewApp("test")
	err := app.validateFlags(
		"",
		"",
		"input.txt",
		[]string{"ds.yaml"},
		nil,
		nil,
		"",
	)
	require.NoError(t, err)
}

func TestValidateFlagsNoData(t *testing.T) {
	app := NewApp("test")

	err := app.validateFlags(
		"",
		"",
		"input.txt",
		nil,
		nil,
		nil,
		"",
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrDataRequired)
}

func TestValidateFlagsNoDataWithEnvsubstEngine(t *testing.T) {
	app := NewApp("test")

	err := app.validateFlags(
		"",
		"",
		"input.txt",
		nil,
		nil,
		nil,
		"envsubst",
	)
	require.NoError(t, err)
}

func TestValidateFlagsNoInput(t *testing.T) {
	app := NewApp("test")
	err := app.validateFlags(
		"",
		"",
		"",
		[]string{"ds.yaml"},
		nil,
		nil,
		"",
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNoInput)
}

func TestValidateFlagsInputFileAndDirConflict(t *testing.T) {
	app := NewApp("test")
	err := app.validateFlags(
		"",
		"input/",
		"input.txt",
		[]string{"ds.yaml"},
		nil,
		nil,
		"",
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInputFileAndDirConflict)
}

func TestValidateFlagsInputStringAndFileConflict(t *testing.T) {
	app := NewApp("test")
	err := app.validateFlags(
		"input-string",
		"",
		"input",
		[]string{"ds.yaml"},
		nil,
		nil,
		"",
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInputStringAndFileConflict)
}

func TestValidateFlagsInputStringAndDirConflict(t *testing.T) {
	app := NewApp("test")
	err := app.validateFlags(
		"input-string",
		"input/",
		"",
		[]string{"ds.yaml"},
		nil,
		nil,
		"",
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInputStringAndDirConflict)
}

func TestValidateFlagsInputFileAndExcludeConflict(t *testing.T) {
	app := NewApp("test")
	err := app.validateFlags(
		"",
		"",
		"input.txt",
		[]string{"ds.yaml"},
		nil,
		[]string{"exclude.txt"},
		"",
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInputFileAndExcludeConflict)
}

func TestValidateFlagsInputStringAndExcludeConflict(t *testing.T) {
	app := NewApp("test")
	err := app.validateFlags(
		"input-string",
		"",
		"",
		[]string{"ds.yaml"},
		nil,
		[]string{"exclude.txt"},
		"",
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInputStringAndExcludeConflict)
}
