package initialize

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/ldelossa/goblog"
	git "github.com/ldelossa/goblog/git"
	"github.com/ldelossa/goblog/pkg/golog"
)

// Upgrade is a Decision which determines if
// a GoBlog upgrade should take place.
//
// If an upgrade should take place it's Yes branch is called.
// If it should not its No branch is called.
func Upgrade(ctx context.Context) (bool, error) {
	action := fmt.Sprintf("A new GoBlog version is available. Use 'goblog upgrade' to upgrade")
	availableVersion, err := git.GlobalGit.LatestTagOrCommit(ctx)
	if err != nil {
		return false, err
	}

	// can we parse it as a semver?
	availableSem, err := semver.NewVersion(strings.TrimPrefix(availableVersion, "v"))
	if err != nil {
		return false, nil
	}

	// can we parse current version as semver?
	thisVersion, err := semver.NewVersion(strings.TrimPrefix(goblog.Version, "v"))
	if err != nil {
		// if no, inform the user there's an upgrade anyway
		// we are probably running a build of a specific commit.
		golog.Warning(action)
		return true, nil
	}

	// if yes, check if this version is less then available
	if thisVersion.LessThan(*availableSem) {
		golog.Warning(action)
		return true, nil
	}
	return false, nil
}
