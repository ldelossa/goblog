package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/fatih/color"
	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/cmd/goblog/internal/config"
	"github.com/ldelossa/goblog/cmd/goblog/internal/diff"
	"github.com/ldelossa/goblog/cmd/goblog/internal/drafts"
	"github.com/ldelossa/goblog/cmd/goblog/internal/initialize"
	"github.com/ldelossa/goblog/cmd/goblog/internal/posts"
	"github.com/ldelossa/goblog/cmd/goblog/internal/serve"
	"github.com/ldelossa/goblog/cmd/goblog/internal/upgrade"
	"github.com/ldelossa/goblog/pkg/golog"
)

const usage = `The goblog command line serves two purposes.
First it may act as an http server, serving assets and blog posts.
Secondly it helps you write and format blog posts. 

The command is split into subcommands, each containing their own help content.

goblog init    - create a new goblog environment
goblog config  - update configuration details
goblog serve   - serve your blog posts, assests, and web root over http
goblog posts   - list, view, and remove published posts
goblog drafts  - list, create, publish, and delete draft blog posts
goblog build   - build a new goblog binary with the latest posts and web root
goblog diff    - diff the contents of your local and embedded tree
goblog preview - preview your blog by running the code in $HOME/src directly
goblog upgrade - upgrade goblog to the newest or specific version
`

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Error: subcommand required\n\n")
		fmt.Println(usage)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	switch os.Args[1] {
	case "--help":
		fmt.Println(usage)
		os.Exit(0)
	case "--version":
		fmt.Println(goblog.Version)
		os.Exit(0)
	case "init":
		initialize.Initialize(ctx)
		os.Exit(0)
	case "config":
		initialize.Initialize(ctx)
		config.Root(ctx)
	case "serve":
		serve.Serve()
	case "posts":
		initialize.Initialize(ctx)
		posts.Root(ctx)
	case "drafts":
		initialize.Initialize(ctx)
		drafts.Root(ctx)
	case "build":
		_, err := initialize.Build(ctx)
		if err != nil {
			initialize.Initialize(ctx)
		}
	case "preview":
		_, err := initialize.Build(ctx)
		cmd := exec.Command("go", "run", "./cmd/goblog/", "serve")
		cmd.Dir = goblog.Src
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			color.Red("Error: failed to start goblog from src directory: %v", err)
			os.Exit(1)
		}
		os.Exit(0)
	case "diff":
		err := diff.Diff(ctx)
		if err != nil {
			golog.Error("%v", err)
			os.Exit(1)
		}
	case "upgrade":
		err := upgrade.Upgrade(ctx)
		if err != nil {
			golog.Error("%v", err)
			os.Exit(1)
		}
	default:
		fmt.Println(usage)
		fmt.Printf("Error: unrecognized subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}
}
