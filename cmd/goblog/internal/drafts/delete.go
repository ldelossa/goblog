package drafts

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/pkg/golog"
)

var deleteFS = flag.NewFlagSet("delete", flag.ExitOnError)

var deleteFlags = struct {
}{}

func delete(ctx context.Context) error {
	deleteFS.Usage = func() {
		fmt.Printf(`
The delete subcommand removes a draft.

Usage:
	goblog drafts delete ID

`)
	}

	if len(os.Args) < 4 {
		editFS.Usage()
		return fmt.Errorf("Error: Not enough arguments provided to 'delete' subcommand\n")
	}

	// first arg must be id
	id, err := strconv.Atoi(os.Args[3])
	if err != nil {
		editFS.Usage()
		return fmt.Errorf("Error: first argument to 'delete' subcommand must be an integer id")
	}

	sorted := goblog.LocalDraftsCache
	if err != nil {
		return fmt.Errorf("Error: failed retrieving drafts: %v", err)
	}
	if len(sorted) == 0 {
		golog.Info(`There are no drafts to delete currently.

Use 'goblog drafts new' to create one.`)
	}

	if id == 0 {
		return fmt.Errorf("Error: must provide a draft id.")
	}

	if id > len(sorted) {
		return fmt.Errorf("Error: draft id %d does not exist", id)
	}

	draft := sorted[id-1]

	err = os.Remove(filepath.Join(goblog.Src, draft.Path))
	if err != nil {
		return fmt.Errorf("Error: failed to remove your draft: %v", err)
	}
	golog.Info(`Successfully deleted draft %v`, id)
	return nil
}
