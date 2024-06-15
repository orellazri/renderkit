package app

import (
	"errors"
)

var (
	ErrNoInput                       = errors.New("input is required")
	ErrInputStringAndDirConflict     = errors.New("only one of input or input-dir can be set")
	ErrInputStringAndFileConflict    = errors.New("only one of input or file can be set")
	ErrInputFileAndDirConflict       = errors.New("only one of input or file can be set")
	ErrInputFileAndExcludeConflict   = errors.New("exclude cannot be used with file")
	ErrInputStringAndExcludeConflict = errors.New("exclude cannot be used with input string")
	ErrNoOutput                      = errors.New("output is required")
	ErrDataRequired                  = errors.New("data is required through the datasource or data flags")
)

func (a *App) validateFlags(
	inputString string,
	inputDir string,
	inputFile string,
	datasource []string,
	data []string,
	excluded []string,
) error {
	if len(inputString) == 0 && len(inputDir) == 0 && len(inputFile) == 0 {
		return ErrNoInput
	}

	if len(inputString) > 0 && len(inputDir) > 0 {
		return ErrInputStringAndDirConflict
	}

	if len(inputString) > 0 && len(inputFile) > 0 {
		return ErrInputStringAndFileConflict
	}

	if len(inputDir) > 0 && len(inputFile) > 0 {
		return ErrInputFileAndDirConflict
	}

	if len(datasource) == 0 && len(data) == 0 {
		return ErrDataRequired
	}

	if len(inputFile) > 0 && len(excluded) > 0 {
		return ErrInputFileAndExcludeConflict
	}

	if len(inputString) > 0 && len(excluded) > 0 {
		return ErrInputStringAndExcludeConflict
	}

	return nil
}
