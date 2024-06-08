package app

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func validateFlags(cCtx *cli.Context) error {
	if len(cCtx.StringSlice("input-file")) != 0 && len(cCtx.String("input-dir")) != 0 {
		return fmt.Errorf("the flags \"input-file\" and \"input-dir\" are mutually exclusive")
	}

	if len(cCtx.String("output-file")) != 0 && len(cCtx.String("output-dir")) != 0 {
		return fmt.Errorf("the flags \"output-file\" and \"output-dir\" are mutually exclusive")
	}

	if len(cCtx.String("input-dir")) != 0 && len(cCtx.String("output-dir")) == 0 {
		return fmt.Errorf("if input-dir is present, output-dir must be present")
	}

	if len(cCtx.StringSlice("input-file")) == 1 && cCtx.String("output-dir") != "" {
		return fmt.Errorf("if input-file has one file, output-dir must not be present")
	}

	if len(cCtx.StringSlice("input-file")) > 1 && cCtx.String("output-file") != "" {
		return fmt.Errorf("if multiple input-files are present, output-file must not be present")
	}

	if len(cCtx.StringSlice("datasource")) == 0 {
		return fmt.Errorf("datasource is required")
	}

	if len(cCtx.String("engine")) == 0 {
		return fmt.Errorf("engine is required")
	}

	return nil
}
