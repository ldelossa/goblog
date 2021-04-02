package posts

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/pkg/golog"
	"gopkg.in/yaml.v3"
)

var viewFS = flag.NewFlagSet("view", flag.ExitOnError)

var viewFlags = struct {
}{}

func view(ctx context.Context) error {
	viewFS.Usage = func() {
		fmt.Printf(`
The view subcommand prints post contents to stdout. 

The '--meta' flag may be used to print both the post content and the metdata data in yaml syntax.

Usage:
	goblog posts view ID [--meta]

`)
	}

	if len(os.Args) < 4 {
		listFS.Usage()
		return fmt.Errorf("Error: Not enough arguments to 'view' subcommand")
	}

	// first arg must be id
	id, err := strconv.Atoi(os.Args[3])
	if err != nil {
		listFS.Usage()
		return fmt.Errorf("Error: first argument to 'view' subcommand must be an integer id")
	}

	var local bool
	for _, arg := range os.Args {
		if arg == "--local" || arg == "-local" {
			local = true
		}
	}

	var meta bool
	for _, arg := range os.Args {
		if arg == "--meta" || arg == "-meta" {
			meta = true
		}
	}

	var posts goblog.DateSortable
	if local {
		posts = goblog.LocalPostsCache
		if err != nil {
			return fmt.Errorf("Error: failed to get local posts: %v", err)
		}
	} else {
		posts = goblog.EmbeddedPostsCache
	}

	if len(posts) == 0 {
		golog.Info(`There are no posts to view currently.

Use 'goblog drafts new' to create one and 'goblog publish' to build a GoBlog binary with your new posts.`)
	}

	if id == 0 {
		return fmt.Errorf("Error: must provide a post id.")
	}

	if id > len(posts) {
		return fmt.Errorf("id not found")
	}

	post := posts[id-1]

	var f fs.File
	if local {
		f, err = os.Open(post.Path)
		if err != nil {
			return fmt.Errorf("error viewing post: " + err.Error())
		}
	} else {
		f, err = goblog.DraftsFS.Open(post.Path)
		if err != nil {
			return fmt.Errorf("error viewing post: " + err.Error())
		}
	}
	// just write out the file data and exit 0
	if meta {
		_, err := io.Copy(os.Stdout, f)
		if err != nil {
			return fmt.Errorf("error viewing post: " + err.Error())
		}
	}

	err = yaml.NewDecoder(f).Decode(&post)
	if err != nil {
		return fmt.Errorf("error viewing post: " + err.Error())
	}

	if meta {
		fmt.Println(post.MarkDown.Value)
	}
	fmt.Println(post.MarkDown.Value)

	return nil
}
