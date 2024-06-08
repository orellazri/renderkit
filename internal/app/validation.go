package app

import (
	"errors"
	"fmt"

	"github.com/orellazri/renderkit/internal/engine"
	"github.com/urfave/cli/v2"
)

func (a *App) validateFlags(cCtx *cli.Context) error {
	if len(cCtx.StringSlice("input")) > 0 && len(cCtx.String("input-dir")) > 0 {
		return fmt.Errorf("the flags \"input\" and \"input-dir\" are mutually exclusive")
	}

	if len(cCtx.String("output")) > 0 && len(cCtx.String("output-dir")) > 0 {
		return fmt.Errorf("the flags \"output\" and \"output-dir\" are mutually exclusive")
	}

	if len(cCtx.String("input-dir")) > 0 && len(cCtx.String("output-dir")) == 0 {
		return fmt.Errorf("if input-dir is present, output-dir must be present")
	}

	if len(cCtx.StringSlice("input")) == 1 && len(cCtx.String("output-dir")) > 0 {
		return fmt.Errorf("if input has one file, output-dir must not be present")
	}

	if len(cCtx.StringSlice("input")) > 1 && len(cCtx.String("output")) > 0 {
		return fmt.Errorf("if multiple inputs are present, output must not be present")
	}

	if len(cCtx.StringSlice("datasource")) == 0 {
		return fmt.Errorf("datasource is required")
	}

	if len(cCtx.String("engine")) == 0 {
		return fmt.Errorf("engine is required")
	}

	return nil
}

func (a *App) setMode(input []string, inputDir string) error {
	if len(input) == 1 {
		a.mode = ModeFileToFile
	} else if len(input) > 1 {
		a.mode = ModeFilesToDir
	} else if len(inputDir) > 0 {
		a.mode = ModeDirToDir
	} else {
		return errors.New("unsupported mode")
	}

	return nil
}

func (a *App) setEngine(engineStr string) error {
	eng, ok := engine.EnginesMap[engineStr]
	if !ok {
		return fmt.Errorf("engine %s not found", engineStr)
	}
	a.engine = eng

	return nil
}
