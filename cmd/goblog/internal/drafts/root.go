package drafts

import (
	"context"
	"fmt"
	"os"

	"github.com/ldelossa/goblog/pkg/golog"
)

var usage = `The 'drafts' subcommand is used to edit and publish local drafts.

Drafts are your working space where you can leave ideas hanging around. 

Once a draft is finished it is 'published' and will be embedded into the next GoBlog binary you build via the 'goblog publish' command.

If you are looking to work with published posts see the 'goblog posts' command instead.

goblog drafts new     - create and edit a new draft
goblog drafts edit    - edit an existing draft or its metadata
goblog drafts list    - list drafts 
goblog drafts view    - view the contents of a draft
goblog drafts delete  - delete a draft
goblog drafts publish - publishes a draft 
`

// Root is the 'drafts' subcommand root handler.
func Root(ctx context.Context) {
	if len(os.Args) < 3 {
		golog.Error("Error: The 'drafts' subcommand requires a further directive.")
		golog.Info(usage)
		os.Exit(1)
	}
	if os.Args[2] == "--help" || os.Args[2] == "-help" {
		fmt.Printf(usage)
		os.Exit(0)
	}
    var err error
	switch os.Args[2] {
	case "edit":
        err = edit(ctx)
	case "new":
        err = new(ctx)
	case "list":
        err = list(ctx)
	case "view":
        err = view(ctx)
	case "delete":
        err = delete(ctx)
	case "publish":
        err = publish(ctx)
	default:
		golog.Error(`Error: unknown subcommand provided.`)
		fmt.Printf(usage)
		os.Exit(1)
	}
    if err != nil {
        golog.Error("%v", err)
        os.Exit(1)
    }
}
