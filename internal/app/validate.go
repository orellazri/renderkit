package app

import "errors"

var (
	ErrNoInput                  = errors.New("either input or input-dir is required")
	ErrNoOuput                  = errors.New("either output or output-dir is required")
	ErrInputExclusive           = errors.New("the flags \"input\" and \"input-dir\" are mutually exclusive")
	ErrOutputExclusive          = errors.New("the flags \"output\" and \"output-dir\" are mutually exclusive")
	ErrInputDirWithoutOutputDir = errors.New("if input-dir is present, output-dir must be present")
	ErrInputFileWithOutputDir   = errors.New("if input has one file, output-dir must not be present")
	ErrMultipleInputsWithOutput = errors.New("if multiple inputs are present, output must not be present")
	ErrDatasourceRequired       = errors.New("datasource is required")
	ErrEngineRequired           = errors.New("engine is required")
)

func (a *App) validateFlags(
	input []string,
	inputDir string,
	output string,
	outputDir string,
	datasource []string,
	engine string,
) error {
	if len(input) == 0 && len(inputDir) == 0 {
		return ErrNoInput
	}

	if len(output) == 0 && len(outputDir) == 0 {
		return ErrNoOuput
	}

	if len(input) > 0 && len(inputDir) > 0 {
		return ErrInputExclusive
	}

	if len(output) > 0 && len(outputDir) > 0 {
		return ErrOutputExclusive
	}

	if len(inputDir) > 0 && len(outputDir) == 0 {
		return ErrInputDirWithoutOutputDir
	}

	if len(input) == 1 && len(outputDir) > 0 {
		return ErrInputFileWithOutputDir
	}

	if len(input) > 1 && len(output) > 0 {
		return ErrMultipleInputsWithOutput
	}

	if len(datasource) == 0 {
		return ErrDatasourceRequired
	}

	if len(engine) == 0 {
		return ErrEngineRequired
	}

	return nil
}
