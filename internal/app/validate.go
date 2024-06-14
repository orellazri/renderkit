package app

import "errors"

var (
	ErrNoInput                 = errors.New("input is required")
	ErrInputFileAndDirConflict = errors.New("only one of input or file can be set")
	ErrNoOuput                 = errors.New("output is required")
	ErrDataRequired            = errors.New("data is required through the datasource or data flags")
)

func (a *App) validateFlags(
	inputDir string,
	inputFile string,
	output string,
	datasource []string,
	data []string,
) error {
	if len(inputDir) == 0 && len(inputFile) == 0 {
		return ErrNoInput
	}

	if len(inputDir) > 0 && len(inputFile) > 0 {
		return ErrInputFileAndDirConflict
	}

	if len(output) == 0 {
		return ErrNoOuput
	}

	if len(datasource) == 0 && len(data) == 0 {
		return ErrDataRequired
	}

	return nil
}
