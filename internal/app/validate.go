package app

import (
	"fmt"
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
		return fmt.Errorf("either input or input-dir is required")
	}

	if len(output) == 0 && len(outputDir) == 0 {
		return fmt.Errorf("either output or output-dir is required")
	}

	if len(input) > 0 && len(inputDir) > 0 {
		return fmt.Errorf("the flags \"input\" and \"input-dir\" are mutually exclusive")
	}

	if len(output) > 0 && len(outputDir) > 0 {
		return fmt.Errorf("the flags \"output\" and \"output-dir\" are mutually exclusive")
	}

	if len(inputDir) > 0 && len(outputDir) == 0 {
		return fmt.Errorf("if input-dir is present, output-dir must be present")
	}

	if len(input) == 1 && len(outputDir) > 0 {
		return fmt.Errorf("if input has one file, output-dir must not be present")
	}

	if len(input) > 1 && len(output) > 0 {
		return fmt.Errorf("if multiple inputs are present, output must not be present")
	}

	if len(datasource) == 0 {
		return fmt.Errorf("datasource is required")
	}

	if len(engine) == 0 {
		return fmt.Errorf("engine is required")
	}

	return nil
}
