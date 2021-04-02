package drafts

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ldelossa/goblog"
)

var listFS = flag.NewFlagSet("list", flag.ExitOnError)

var listFlags = struct {
}{}

func list(ctx context.Context) error {
	listFS.Usage = func() {
		fmt.Printf(`
The list subcommand lists drafts and their ids. 

If the "--embedded" flag is provided only drafts embedded into the current GoBlog binary will be displayed.

This subcommand takes no arguments.

Usage:
	goblog drafts --embedded list

`)
	}

	for _, arg := range os.Args {
		if arg == "--help" || arg == "-help" {
			listFS.Usage()
			return nil
		}
	}

	var embedded bool
	for _, arg := range os.Args {
		if arg == "--embedded" || arg == "-embedded" {
			embedded = true
		}
	}

	var sorted goblog.DateSortable
	var err error
	if embedded {
		sorted = goblog.EmbeddedDraftsCache
	} else {
		sorted = goblog.LocalDraftsCache
		if err != nil {
			return fmt.Errorf("Error: failed retrieving drafts: %v", err)
		}
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	if embedded {
		fmt.Fprintln(tw, "ID\tDATE\tTITLE\tSUMMARY\t(embedded)")
	} else {
		fmt.Fprintln(tw, "ID\tDATE\tTITLE\tSUMMARY\t(local)")
	}

	for i, draft := range sorted {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n", i+1, draft.Date.Format("2006-Jan-2"), draft.Title, draft.Summary)
	}
	err = tw.Flush()
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	return nil
}
