package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateFlagsNoErrors(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		"input.txt",
		"output/",
		[]string{"ds.yaml"},
		"gotemplates",
	)
	require.NoError(t, err)
}

func TestValidateFlagsNoDatasource(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		"input.txt",
		"output/",
		nil,
		"gotemplates",
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrDatasourceRequired)
}

func TestValidateFlagsNoEngine(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		"input.txt",
		"output/",
		[]string{"ds.yaml"},
		"",
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrEngineRequired)
}

func TestValidateFlagsNoInput(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		"",
		"output/",
		[]string{"ds.yaml"},
		"gotemplates",
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNoInput)
}

func TestValidateFlagsNoOutput(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		"input.txt",
		"",
		[]string{"ds.yaml"},
		"gotemplates",
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNoOuput)
}
