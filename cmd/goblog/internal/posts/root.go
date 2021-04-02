package posts

import (
	"context"
	"os"

	"github.com/ldelossa/goblog/pkg/golog"
)

var usage = `The 'posts' subcommand is for managing published blog posts.
These posts are embedded into the GoBlog binary.
If you're removing a post you'll need to rebuild GoBlog.

The '--local' flag optionally instructs GoBlog is look at local posts, ones not
embedded into the binary.

Usage: 

goblog posts list  - list published blog posts and their id
goblog posts view  - view the markdown contents of a post
goblog posts draft - unpublish a post and move it to draft (assumes --local flag)

`

// Root is the 'posts' subcommand root handler
func Root(ctx context.Context) {
	if len(os.Args) < 3 {
		golog.Info(usage)
		golog.Error(`Error: The 'posts' subcommand requires a directive.`)
		os.Exit(1)
	}

    var err error
	switch os.Args[2] {
	case "--help":
		golog.Info(usage)
		os.Exit(0)
	case "list":
        err = list(ctx)
	case "draft":
        err = draft(ctx)
	case "view":
		err = view(ctx)
	default:
		golog.Fatal(`Error: unrecognized subcommand: %s`, os.Args[2])
	}
    if err != nil {
        golog.Error("%v", err)
        os.Exit(1)
    }
}
