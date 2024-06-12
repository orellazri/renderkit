package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateFlagsNoErrors(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		[]string{"file.txt"},
		"",
		"output.txt",
		"",
		[]string{"ds.yaml"},
		"gotemplates",
	)
	require.NoError(t, err)
}

func TestValidateFlagsNoDatasource(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		[]string{"file.txt"},
		"",
		"output.txt",
		"",
		nil,
		"gotemplates",
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "datasource is required")
}

func TestValidateFlagsNoEngine(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		[]string{"file.txt"},
		"",
		"output.txt",
		"",
		[]string{"ds.yaml"},
		"",
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "engine is required")
}

func TestValidateFlagsNoInput(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		[]string{},
		"",
		"output.txt",
		"",
		[]string{"ds.yaml"},
		"gotemplates",
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "either input or input-dir is required")
}

func TestValidateFlagsNoOutput(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		[]string{"file.txt"},
		"",
		"",
		"",
		[]string{"ds.yaml"},
		"gotemplates",
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "either output or output-dir is required")
}

func TestValidateFlagsInputAndInputDir(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		[]string{"file.txt"},
		"dir",
		"output.txt",
		"",
		[]string{"ds.yaml"},
		"gotemplates",
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "the flags \"input\" and \"input-dir\" are mutually exclusive")
}

func TestValidateFlagsOutputAndOutputDir(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		[]string{"file.txt"},
		"",
		"output.txt",
		"dir",
		[]string{"ds.yaml"},
		"gotemplates",
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "the flags \"output\" and \"output-dir\" are mutually exclusive")
}

func TestValidateFlagsInputDirAndOutputDir(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		[]string{},
		"dir",
		"out.txt",
		"",
		[]string{"ds.yaml"},
		"gotemplates",
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "if input-dir is present, output-dir must be present")
}

func TestValidateFlagsInputAndOutputDir(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		[]string{"file.txt"},
		"",
		"",
		"dir",
		[]string{"ds.yaml"},
		"gotemplates",
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "if input has one file, output-dir must not be present")
}

func TestValidateFlagsMultipleInputsAndOutput(t *testing.T) {
	app := NewApp()
	err := app.validateFlags(
		[]string{"file1.txt", "file2.txt"},
		"",
		"output.txt",
		"",
		[]string{"ds.yaml"},
		"gotemplates",
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "if multiple inputs are present, output must not be present")
}
