package app

import (
	"fmt"
	"os"

	"github.com/orellazri/renderkit/internal/engine"
	"github.com/urfave/cli/v2"
)

func renderTemplate(cCtx *cli.Context, inFile *os.File, outFile *os.File, data map[string]any) error {
	eng, ok := engine.EnginesMap[cCtx.String("engine")]
	if !ok {
		return fmt.Errorf("engine %s not found", cCtx.String("engine"))
	}
	if err := eng.Render(inFile.Name(), outFile, data); err != nil {
		return fmt.Errorf("render: %s", err)
	}

	return nil
}
