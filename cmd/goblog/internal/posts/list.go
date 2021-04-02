package posts

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
The 'list' subcommand will list posts in date order.

If the '--local' flag is used a list of local posts, ones not emedded into the binary, will be listed.

This subcommand takes no arguments.
`)
	}

	var local bool
	for _, arg := range os.Args {
		if arg == "--local" || arg == "-local" {
			local = true
		}
	}

	var posts goblog.DateSortable
	var err error
	if local {
		posts = goblog.LocalPostsCache
		if err != nil {
			return fmt.Errorf("Error: failed to query local posts: %v", err)
		}
	} else {
		posts = goblog.EmbeddedPostsCache
	}

	if len(posts) == 0 {
		fmt.Println("No posts found.")
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	if local {
		fmt.Fprintln(tw, "ID\tDATE\tTITLE\tSUMMARY\t(local)")
	} else {
		fmt.Fprintln(tw, "ID\tDATE\tTITLE\tSUMMARY\t(embedded)")
	}
	for i, post := range posts {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n", i+1, post.Date.Format("2006-Jan-2"), post.Title, post.Summary)
	}
	err = tw.Flush()
	if err != nil {
		return fmt.Errorf("error: " + err.Error())
	}
	return nil
}
