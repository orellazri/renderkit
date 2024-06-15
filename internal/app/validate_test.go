package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateFlagsNoErrors(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		"",
		"",
		"input.txt",
		[]string{"ds.yaml"},
		nil,
	)
	require.NoError(t, err)
}

func TestValidateFlagsNoData(t *testing.T) {
	app := NewApp()

	err := app.validateFlags(
		"",
		"",
		"input.txt",
		nil,
		nil,
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrDataRequired)
}

func TestValidateFlagsNoInput(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		"",
		"",
		"",
		[]string{"ds.yaml"},
		nil,
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNoInput)
}

func TestValidateFlagsInputFileAndDirConflict(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		"",
		"input/",
		"input.txt",
		[]string{"ds.yaml"},
		nil,
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInputFileAndDirConflict)
}

func TestValidateFlagsInputStringAndFileConflict(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		"input-string",
		"",
		"input",
		[]string{"ds.yaml"},
		nil,
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInputStringAndFileConflict)
}

func TestValidateFlagsInputStringAndDirConflict(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		"input-string",
		"input/",
		"",
		[]string{"ds.yaml"},
		nil,
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInputStringAndDirConflict)
}
