package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateFlagsNoErrors(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		"",
		"input.txt",
		"output/",
		[]string{"ds.yaml"},
		nil,
	)
	require.NoError(t, err)
}

func TestValidateFlagsNoData(t *testing.T) {
	app := NewApp()

	err := app.validateFlags(
		"",
		"input.txt",
		"output/",
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
		"output/",
		[]string{"ds.yaml"},
		nil,
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNoInput)
}

func TestValidateFlagsInputFileAndDirConflict(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		"input/",
		"input.txt",
		"output/",
		[]string{"ds.yaml"},
		nil,
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInputFileAndDirConflict)
}

func TestValidateFlagsNoOutput(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		"",
		"input.txt",
		"",
		[]string{"ds.yaml"},
		nil,
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNoOuput)
}
