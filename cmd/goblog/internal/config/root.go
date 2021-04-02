package config

import (
	"context"
	"fmt"
	"os"

	"github.com/ldelossa/goblog/cmd/goblog/internal/initialize"
	"github.com/ldelossa/goblog/pkg/golog"
)

var usage = `The 'config' subcommand is for updating configuration.
These configs are embedded into the GoBlog binary.
If you're changing a config option you'll need to rebuild GoBlog.

goblog config app-paths  - specify your web applicatoin's
goblog config fork       - update your goblog fork
`

func Root(ctx context.Context) {
	if len(os.Args) < 3 {
		golog.Error("Error: The 'config' subcommand requires a directive.\n")
		golog.Info(usage)
		os.Exit(1)
	}
	if os.Args[2] == "--help" || os.Args[2] == "-help" {
		fmt.Printf(usage)
		os.Exit(0)
	}
	initialize.Initialize(context.TODO())
	switch os.Args[2] {
	case "app-paths":
		appPaths(ctx)
	case "fork":
	}
}
