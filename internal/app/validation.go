package app

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func checkFlag(cCtx *cli.Context, flag1, flag2 string) error {
	if len(cCtx.String(flag1)) != 0 && len(cCtx.String(flag2)) != 0 {
		return fmt.Errorf("the flags \"%s\" and \"%s\" are mutually exclusive", flag1, flag2)
	}
	if len(cCtx.String(flag1)) == 0 && len(cCtx.String(flag2)) == 0 {
		return fmt.Errorf("either the \"%s\" flag or the \"%s\" flag is required", flag1, flag2)
	}
	return nil
}

func validateFlags(cCtx *cli.Context) error {
	if err := checkFlag(cCtx, "input", "input-dir"); err != nil {
		return err
	}

	if err := checkFlag(cCtx, "output", "output-dir"); err != nil {
		return err
	}

	if len(cCtx.String("input-dir")) != 0 && len(cCtx.String("output-dir")) == 0 {
		return fmt.Errorf("if input-dir is present, output-dir must be present")
	}

	if len(cCtx.String("output-dir")) != 0 && len(cCtx.String("input-dir")) == 0 {
		return fmt.Errorf("if output-dir is present, input-dir must be present")
	}

	if len(cCtx.StringSlice("datasource")) == 0 {
		return fmt.Errorf("datasource is required")
	}

	if len(cCtx.String("engine")) == 0 {
		return fmt.Errorf("engine is required")
	}

	return nil
}
