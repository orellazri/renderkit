package app

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func validateFlags(cCtx *cli.Context) error {
	if len(cCtx.String("input")) == 0 {
		return fmt.Errorf("input file is required")
	}

	if len(cCtx.StringSlice("datasource")) == 0 {
		return fmt.Errorf("datasource is required")
	}

	if len(cCtx.String("engine")) == 0 {
		return fmt.Errorf("engine is required")
	}

	if len(cCtx.String("output")) == 0 {
		return fmt.Errorf("output file is required")
	}

	return nil
}
