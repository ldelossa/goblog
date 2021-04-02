package posts

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

var draftFS = flag.NewFlagSet("draft", flag.ExitOnError)

var draftFlags = struct {
}{}

func draft(ctx context.Context) error {
	listFS.Usage = func() {
		fmt.Printf(`
The 'draft' subcommand moves published posts into drafts.

If the '--local' flag is used it will target only local posts not embedded in a GoBlog binary yet.

Unless the '--local' flag was used, you will need to issue a 'goblog publish' to see the results in a new GoBlog binary.

Usage:
	goblog posts draft ID

`)
	}

	// 0: goblog 1: posts 2: draft
	if len(os.Args) < 4 {
		listFS.Usage()
		return fmt.Errorf(`Error: Not enough arguments to 'draft' subcommand`)
	}

	id, err := strconv.Atoi(os.Args[3])
	if err != nil {
		listFS.Usage()
		return fmt.Errorf(`Error: 'draft' subcommand requires an integer ID argument`)
	}

	var local bool
	for _, arg := range os.Args {
		if arg == "--local" || arg == "-local" {
			local = true
		}
	}

	var posts goblog.DateSortable
	if local {
		posts = goblog.LocalPostsCache
		if err != nil {
			return fmt.Errorf(`Error: failed to get local posts: %v`, err)
		}
	} else {
		posts = goblog.EmbeddedPostsCache
	}

	if len(posts) == 0 {
		golog.Info(`There are no posts to view currently.\nUse 'goblog drafts new' to create one and 'goblog build' to build a GoBlog binary with your new posts.`)
	}

	if id == 0 {
		return fmt.Errorf(`Error: must provide a post id.`)
	}

	if id > len(posts) {
		return fmt.Errorf(`id not found`)
	}

	post := posts[id-1]
	base := filepath.Base(post.Path)
	err = os.Rename(path.Join(goblog.Posts, base), path.Join(goblog.Drafts, base))
	if err != nil {
		return fmt.Errorf(`Error: failed to move post from posts dir to drafts dir: %v`, err)
	}

	golog.Info(`Post %v successfully moved to drafts.

Use 'goblog build' to build a GoBlog binary with this post removed.`, id)

	return nil
}
