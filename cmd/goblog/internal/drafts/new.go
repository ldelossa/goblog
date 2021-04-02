package drafts

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/cmd/goblog/internal/initialize"
	"github.com/ldelossa/goblog/pkg/golog"
	"github.com/ldelossa/goblog/ui"
	"gopkg.in/yaml.v3"
)

var newFS = flag.NewFlagSet("new", flag.ExitOnError)

var newFlags = struct {
}{}

func new(ctx context.Context) error {
	newFS.Usage = func() {
		fmt.Printf(`
The new subcommand creates a new draft and opens your $EDITOR to it. 

You must save the contents before closing your $EDITOR for GoBlog to correctly save the draft contents.

On close of the $EDITOR you will choose to either publish or leave the draft for later editing.

This subcommand takes no arguments.

Usage:
	goblog drafts new

`)
	}
	// 0: goblog, 1: posts, 2: new
	newFS.Parse(os.Args[3:])

	draft, err := ui.GlobalPrompter.DraftBuilder(ctx)
	if err != nil {
		return fmt.Errorf("Error: failed prompting for post details: %v", err)
	}

	if _, err := os.Stat(goblog.Drafts); os.IsNotExist(err) {
		err := os.Mkdir(goblog.Drafts, 0770)
		if err != nil {
			return fmt.Errorf("Error: failed to create drafts directory: %v", err)
		}
	}

	editor, err := ui.NewEditor(ctx, &draft)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	err = editor.Edit(ctx)
	if err != nil {
		return fmt.Errorf("Error: failed to create drafts directory: %v", err)
	}

	// ask user if they want to publish this or keep draft
	publish, err := ui.GlobalPrompter.ShouldPublishPost(ctx)
	var postPath string
	base := filepath.Base(draft.Path)
	if publish {
		postPath = path.Join(goblog.Posts, base)
	} else {
		postPath = path.Join(goblog.Drafts, base)
	}

	f, err := os.OpenFile(postPath, os.O_CREATE|os.O_WRONLY, 0660)
	onErrDumpMarkdown(draft.MarkDown.Value, err)
	defer f.Close()

	err = yaml.NewEncoder(f).Encode(draft)
	onErrDumpMarkdown(draft.MarkDown.Value, err)

	golog.Info(`Your draft has been written to: %v

It will now be available for usage in subsequent 'drafts' commands.`, postPath)

	// ask user if they want to build a new binary with drafts
	// embedded
	build, err := ui.GlobalPrompter.ShouldBuild(ctx)
	onErrDumpMarkdown(draft.MarkDown.Value, err)

	if build {
		_, err := initialize.Build(ctx)
		if err != nil {
			return fmt.Errorf(`Error: failed to publish new GoBlog binary: %v`, err)
		}
	}

	return nil
}
