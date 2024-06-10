package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateFlagsNoDatasource(t *testing.T) {
	app := NewApp()
	args := []string{"", "--engine", "gotemplates",
		"--input", "file.txt", "--output", "output.txt",
	}
	err := app.Run(args)
	require.Error(t, err)
	require.Contains(t, err.Error(), "datasource is required")
}

func TestValidateFlagsNoEngine(t *testing.T) {
	app := NewApp()
	args := []string{"", "--datasource", "ds.yaml",
		"--input", "file.txt", "--output", "output.txt",
	}
	err := app.Run(args)
	require.Error(t, err)
	require.Contains(t, err.Error(), "engine is required")
}

func TestValidateFlagsInputAndInputDir(t *testing.T) {
	app := NewApp()
	args := []string{"", "--datasource", "ds.yaml", "--engine", "gotemplates",
		"--input", "file.txt", "--input-dir", "dir",
	}
	err := app.Run(args)
	require.Error(t, err)
	require.Contains(t, err.Error(), "the flags \"input\" and \"input-dir\" are mutually exclusive")
}

func TestValidateFlagsOutputAndOutputDir(t *testing.T) {
	app := NewApp()
	args := []string{"", "--datasource", "ds.yaml", "--engine", "gotemplates",
		"--output", "output.txt", "--output-dir", "output_dir",
	}
	err := app.Run(args)
	require.Error(t, err)
	require.Contains(t, err.Error(), "the flags \"output\" and \"output-dir\" are mutually exclusive")
}

func TestValidateFlagsInputDirAndOutputDir(t *testing.T) {
	app := NewApp()
	args := []string{"", "--datasource", "ds.yaml", "--engine", "gotemplates",
		"--input-dir", "input_dir", "--output-dir", "",
	}
	err := app.Run(args)
	require.Error(t, err)
	require.Contains(t, err.Error(), "if input-dir is present, output-dir must be present")
}

func TestValidateFlagsInputAndOutputDir(t *testing.T) {
	app := NewApp()
	args := []string{"", "--datasource", "ds.yaml", "--engine", "gotemplates",
		"--input", "file.txt", "--output-dir", "output_dir",
	}
	err := app.Run(args)
	require.Error(t, err)
	require.Contains(t, err.Error(), "if input has one file, output-dir must not be present")
}

func TestValidateFlagsMultipleInputsAndOutput(t *testing.T) {
	app := NewApp()
	args := []string{"", "--datasource", "ds.yaml", "--engine", "gotemplates",
		"--input", "file1.txt", "--input", "file2.txt", "--output", "output.txt",
	}
	err := app.Run(args)
	require.Error(t, err)
	require.Contains(t, err.Error(), "if multiple inputs are present, output must not be present")
}
