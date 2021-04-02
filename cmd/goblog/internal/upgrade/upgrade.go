package upgrade

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/cmd/goblog/internal/initialize"
	"github.com/ldelossa/goblog/git"
	"github.com/ldelossa/goblog/pkg/golog"
)

var fs = flag.NewFlagSet("serve", flag.ExitOnError)

var flags = struct {
	unstable *bool
	commit   *string
}{
	unstable: fs.Bool("unstable", false, "upgrade GoBlog to the latest upstream commit"),
	commit:   fs.String("commit", "", "upgrade or downgrade GoBlog to the provided commit"),
}

func Upgrade(ctx context.Context) error {
	fs.Usage = func() {
		fmt.Printf(`
The upgrade subcommand is used to build different versions of GoBlog.

By default this subcommand will upgrade GoBlog to the latest upstream tag available. 

If the "--unstable" flag a new GoBlog reflecting the latest upstream commit will be built.

If the "--commit" flag is provided a new GoBlog will be built reflecting the provided commit hash, tag, or branch name.

If utilizing the "--unstable" or "--commit" flags GoBlog will continue to inform you that an upgrade is available pointing
to the latest semver release. 

This is an intentional reminder to run the stable build.

Usage:
	goblog upgrade [--unstable | --commit = COMMIT-ISH]
`)
	}

	// 0: goblog, 1: upgrade
	fs.Parse(os.Args[2:])

	differ := goblog.Differ{}
	// check if there's any diffs, and
	// if so error out.
	draftsDiff, err := differ.Drafts(ctx, false)
	if err != nil {
		return fmt.Errorf(`Error: failed to obtain drafts diff: %v`, err)
	}
	postsDiff, err := differ.Posts(ctx, false)
	if err != nil {
		return fmt.Errorf(`Error: failed to obtain posts diff: %v`, err)
	}
	configDiff, err := differ.Config(ctx)
	if err != nil {
		return fmt.Errorf(`Error: failed to obtain config diff: %v`, err)
	}
	webDiff, err := differ.Web(ctx, false)
	if err != nil {
		return fmt.Errorf(`Error: failed to obtain web diff: %v`, err)
	}

	if draftsDiff != "" || postsDiff != "" || configDiff != "" || webDiff != "" {
		return fmt.Errorf(`Error: 
Local tree and embedded tree are not in sync.
Running "goblog diff" will indicate what is different between trees.`)

	}

	switch {
	case *flags.commit != "":
		err = upgradeCommit(ctx, *flags.commit)
	case *flags.unstable:
		err = upgradeUnstable(ctx)
	default:
		err = upgradeStable(ctx)
	}

	return err

}

func upgradeStable(ctx context.Context) error {
	if err := git.GlobalGit.Reset(ctx); err != nil {
		return err
	}
	if err := git.GlobalGit.Clean(ctx); err != nil {
		return err
	}

	tag, err := git.GlobalGit.LatestTag(ctx)
	if err != nil {
		return fmt.Errorf(`Error: failed to checkout latest tag: %v`, err)
	}

	if tag == "" {
		golog.Warning(`No tag discovered.`)
	}

	_, err = initialize.Synchronize(ctx)
	if err != nil {
		return fmt.Errorf(`Error: failed to synchronize local and embedded tree: %v.`, err)
	}

	_, err = initialize.Build(ctx)
	if err != nil {
		return fmt.Errorf(`Error: failed to upgrade GoBlog: %v.`, err)
	}

	golog.Info(`GoBlog successfully updated to %v`, tag)
	return nil
}

func upgradeUnstable(ctx context.Context) error {
	// just reset to latest commit.
	if err := git.GlobalGit.Clean(ctx); err != nil {
		return err
	}
	if err := git.GlobalGit.Reset(ctx); err != nil {
		return err
	}

	commit, err := git.GlobalGit.HEAD(ctx)
	if err != nil {
		return err
	}

	_, err = initialize.Synchronize(ctx)
	if err != nil {
		return fmt.Errorf(`Error: failed to synchronize local and embedded tree: %v.`, err)
	}

	_, err = initialize.Build(ctx)
	if err != nil {
		return fmt.Errorf(`Error: failed to upgrade GoBlog: %v.`, err)
	}

	golog.Info(`GoBlog successfully updated to %v`, commit)
	return nil
}

// upgradeCommit performs an upgrade/downgrade given
// a specific commit.
func upgradeCommit(ctx context.Context, commit string) error {
	if err := git.GlobalGit.Clean(ctx); err != nil {
		return err
	}
	if err := git.GlobalGit.Reset(ctx); err != nil {
		return err
	}

	if err := git.GlobalGit.Checkout(ctx, *flags.commit); err != nil {
		return err
	}

	_, err := initialize.Synchronize(ctx)
	if err != nil {
		return fmt.Errorf(`Error: failed to synchronize local and embedded tree: %v.`, err)
	}

	_, err = initialize.Build(ctx)
	if err != nil {
		return fmt.Errorf(`Error: failed to upgrade GoBlog: %v.`, err)
	}

	golog.Info(`GoBlog successfully updated to %v`, *flags.commit)
	return nil
}
