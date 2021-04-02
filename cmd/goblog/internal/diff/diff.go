package diff

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/pkg/golog"
)

var fs = flag.NewFlagSet("serve", flag.ExitOnError)

var flags = struct {
	show *bool
}{
	show: fs.Bool("show", false, "show a textual diff along with summary"),
}

func Diff(ctx context.Context) error {
	fs.Usage = func() {
		fmt.Printf(`
The diff subcommand explains any differences between your local and embedded trees. 

A diff line with a "+" indicates a post or draft in your local tree that is not present in your embedded tree. 

A diff line with a "-" indicates a post or draft in your embedded tree not present in your local tree.

If a diff line is suffixed with "dirty" this indicates the post is in both trees but they differ in content. 

You may use the optional "--show" flag to display where dirty posts or drafts differ in their text.

No diff lines indicates your local and embedded trees are in sync.

A call to "goblog publish" will create a new GoBlog binary with any "+" diff lines embedded into the binary.

Usage:
	goblog diff [--show]

`)
	}

	// 0: goblog, 1: diff
	fs.Parse(os.Args[2:])

	differ := goblog.Differ{}

	golog.Info("diff drafts\n- embedded\n+ local\n")
	diff, err := differ.Drafts(context.TODO(), *flags.show)
	if err != nil {
		return fmt.Errorf("Error: failed to get drafts diff: %v", err)
	}
	golog.Warning(diff)

	golog.Info("\ndiff posts\n- embedded\n+ local\n")
	diff, err = differ.Posts(context.TODO(), *flags.show)
	if err != nil {
		return fmt.Errorf("Error: failed to get posts diff: %v", err)
	}
	golog.Warning(diff)

	diff, err = differ.Config(context.TODO())
	if err != nil {
		return fmt.Errorf("Error: failed to get config diff: %v", err)
	}
	if diff != "" {
		golog.Warning("\n- config/config.yaml dirty\n")
	}

	diff, err = differ.Web(context.TODO(), *flags.show)
	if err != nil {
		return fmt.Errorf("Error: failed to get web diff: %v", err)
	}
	if diff != "" {
		golog.Warning("\n- web fs dirty\n")
	}

	return nil
}
