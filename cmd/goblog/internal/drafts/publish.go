package drafts

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/pkg/golog"
)

var publishFS = flag.NewFlagSet("publish", flag.ExitOnError)

var publishFlags = struct{}{}

func publish(ctx context.Context) error {
	publishFS.Usage = func() {
		fmt.Printf(`
The publish subcommand publishes an existing draft. 

When a draft is published it will be embedded into the next GoBlog binary created by running 'goblog publish'.

Usage:
	goblog drafts publish ID

`)
	}

	// 0: goblog, 1: drafts, 2: publish
	publishFS.Parse(os.Args[3:])

	if len(os.Args) < 4 {
		editFS.Usage()
		golog.Error("Error: Not enough arguments provided to 'edit' subcommand\n")
	}

	// first arg must be id
	id, err := strconv.Atoi(os.Args[3])
	if err != nil {
		return fmt.Errorf("Error: first argument to 'publish' subcommand must be an integer id")
	}

	sorted := goblog.LocalDraftsCache
	if err != nil {
		return fmt.Errorf("Error: failed retrieving drafts: %v", err)
	}
	if len(sorted) == 0 {
		golog.Info("There are no drafts to edit currently.\nUse 'goblog drafts new' to create one.")
	}

	if id == 0 {
		return fmt.Errorf("Error: must provide a draft id.")
	}

	if id > len(sorted) {
		return fmt.Errorf("Error: draft id %d does not exist", id)
	}

	draft := sorted[id-1]
	base := filepath.Base(draft.Path)

	err = os.Rename(
		path.Join(goblog.Drafts, base),
		path.Join(goblog.Posts, base),
	)
	if err != nil {
		return fmt.Errorf(`Error: failed to move draft to the posts directory: %v`, err)
	}

	return nil
}
