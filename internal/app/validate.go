package app

import "errors"

var (
	ErrNoInput            = errors.New("input is required")
	ErrNoOuput            = errors.New("output is required")
	ErrDatasourceRequired = errors.New("datasource is required")
	ErrEngineRequired     = errors.New("engine is required")
)

func (a *App) validateFlags(
	input string,
	output string,
	datasource []string,
	engine string,
) error {
	if len(input) == 0 {
		return ErrNoInput
	}

	if len(output) == 0 {
		return ErrNoOuput
	}

	if len(datasource) == 0 {
		return ErrDatasourceRequired
	}

	if len(engine) == 0 {
		return ErrEngineRequired
	}

	return nil
}
