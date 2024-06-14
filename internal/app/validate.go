package app

import "errors"

var (
	ErrNoInput        = errors.New("input is required")
	ErrNoOuput        = errors.New("output is required")
	ErrDataRequired   = errors.New("data is required through the datasource or data flags")
	ErrEngineRequired = errors.New("engine is required")
)

func (a *App) validateFlags(
	input string,
	output string,
	datasource []string,
	data []string,
	engine string,
) error {
	if len(input) == 0 {
		return ErrNoInput
	}

	if len(output) == 0 {
		return ErrNoOuput
	}

	if len(datasource) == 0 && len(data) == 0 {
		return ErrDataRequired
	}

	if len(engine) == 0 {
		return ErrEngineRequired
	}

	return nil
}
