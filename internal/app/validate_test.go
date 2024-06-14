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
		nil,
		"gotemplates",
	)
	require.NoError(t, err)
}

func TestValidateFlagsNoData(t *testing.T) {
	app := NewApp()

	err := app.validateFlags(
		"input.txt",
		"output/",
		nil,
		nil,
		"gotemplates",
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrDataRequired)
}

func TestValidateFlagsNoEngine(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		"input.txt",
		"output/",
		[]string{"ds.yaml"},
		nil,
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
		nil,
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
		nil,
		"gotemplates",
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNoOuput)
}
