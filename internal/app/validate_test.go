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
	require.ErrorIs(t, err, ErrDatasourceRequired)
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
	require.ErrorIs(t, err, ErrEngineRequired)
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
	require.ErrorIs(t, err, ErrNoInput)
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
	require.ErrorIs(t, err, ErrNoOuput)
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
	require.ErrorIs(t, err, ErrInputExclusive)
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
	require.ErrorIs(t, err, ErrOutputExclusive)
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
	require.ErrorIs(t, err, ErrInputDirWithoutOutputDir)
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
	require.ErrorIs(t, err, ErrInputFileWithOutputDir)
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
	require.ErrorIs(t, err, ErrMultipleInputsWithOutput)
}
