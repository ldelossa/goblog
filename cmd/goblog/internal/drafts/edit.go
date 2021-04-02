package drafts

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/cmd/goblog/internal/initialize"
	"github.com/ldelossa/goblog/pkg/golog"
	"github.com/ldelossa/goblog/ui"
	"gopkg.in/yaml.v3"
)

var editFS = flag.NewFlagSet("edit", flag.ExitOnError)

var editFlags = struct {
}{}

func edit(ctx context.Context) error {
	editFS.Usage = func() {
		fmt.Printf(`
The edit subcommand opens a draft for editing. 

The '--meta' flag may be used to edit both the contents and metadata of a post in yaml syntax.

Usage:
	goblog drafts edit ID [--meta]

`)
	}
	// 0: goblog, 1: drafts, 2: edit
	editFS.Parse(os.Args[3:])

	if len(os.Args) < 4 {
		editFS.Usage()
		return fmt.Errorf("Error: Not enough arguments provided to 'edit' subcommand\n")
	}

	// first arg must be id
	id, err := strconv.Atoi(os.Args[3])
	if err != nil {
		editFS.Usage()
		return fmt.Errorf("Error: first argument to 'view' subcommand must be an integer id")
	}

	var meta bool
	for _, arg := range os.Args {
		if arg == "--meta" || arg == "-meta" {
			meta = true
		}
	}

	sorted := goblog.LocalDraftsCache
	if len(sorted) == 0 {
		golog.Info("There are no drafts to edit currently.\nUse 'goblog drafts new' to create one.")
		return nil
	}

	if id == 0 {
		editFS.Usage()
		return fmt.Errorf("Error: must provide a draft id.")
	}

	if id > len(sorted) {
		return fmt.Errorf("Error: draft id %d does not exist", id)
	}

	draft := sorted[id-1]

	editor, err := ui.NewEditor(ctx, &draft)
	if err != nil {
		golog.Fatal("%v", err)
	}

	// handle meta edit only
	if meta {
		// call editor
		err := editor.EditMeta(ctx)
		if err != nil {
			return fmt.Errorf("Error: failed to start editor: %v", err)
		}
		return nil
	}

	err = editor.Edit(ctx)
	if err != nil {
		return fmt.Errorf("Error: failed to start editor: %v", err)
	}

	// ask user if they want to publish this or keep draft
	publish, err := ui.GlobalPrompter.ShouldPublishPost(ctx)
	onErrDumpMarkdown(draft.MarkDown.Value, err)

	// draft.Path will already be formated and have
	// .post syntax
	formated := path.Base(draft.Path)
	var postPath string
	if publish {
		postPath = path.Join(goblog.Posts, formated)
	} else {
		postPath = path.Join(goblog.Drafts, formated)
	}

	f, err := os.OpenFile(postPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	onErrDumpMarkdown(draft.MarkDown.Value, err)
	defer f.Close()

	err = yaml.NewEncoder(f).Encode(draft)
	onErrDumpMarkdown(draft.MarkDown.Value, err)

	golog.Info(`Your draft has been written to: %v.`, postPath)

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

func onErrDumpMarkdown(md string, err error) {
	if err != nil {
		golog.Error("Error: failed editing post: %v", err)
		golog.Error("Dumping your markdown so you don't loose your work...\n")
		fmt.Println("----MARKDOWN BEGIN----")
		fmt.Printf(md)
		fmt.Println("----MARKDOWN END----")
		os.Exit(1)
	}
}
