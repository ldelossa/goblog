package config

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/pkg/golog"
	"gopkg.in/yaml.v3"
)

var appPathsFS = flag.NewFlagSet("app-paths", flag.ExitOnError)

var appPathsFlags = struct{}{}

func appPaths(ctx context.Context) {
	color.Blue(`
Provide a comma separated list of paths your front-end web application will handle.

Each path should have a leading forward slash.

Example: /post,/archive,/settings

`)

	var list string
	_, err := fmt.Scanln(&list)
	if err != nil {
		golog.Fatal("failed to scan input: %v", err)
	}

	paths := strings.Split(list, ",")
	golog.Info("Adding the following paths: %v\n", paths)

	goblog.Conf.AppPaths = paths

	dest := path.Join(goblog.Home, "src/config/config.yaml")
	if _, err := os.Stat(dest); err != nil {
		golog.Fatal("Could not stat config: %v", err)
	}

	f, err := os.OpenFile(dest, os.O_RDWR|os.O_TRUNC, 0)
	if err != nil {
		golog.Fatal("Failed to open config.yaml: %v", err)
	}
	defer f.Close()

	if err = yaml.NewEncoder(f).Encode(&goblog.Conf); err != nil {
		golog.Fatal("Failed to write config: %v", err)
	}
	golog.Info("Wrote new config to %v\n", dest)
}
