package initialize

import (
	"context"

	"github.com/ldelossa/goblog/git"
	"github.com/ldelossa/goblog/pkg/golog"
)

const (
	upstream = "https://github.com/ldelossa/goblog.git"
)

// GitClone clones the upstream GoBlog repository
// into $HOME/src and unpacks any embedded content into this
// working directory.
func GitClone(ctx context.Context) (bool, error) {
	action := "Cloned GoBlog's source code from " + upstream
	err := git.GlobalGit.Clone(ctx, upstream)
	if err != nil {
		return false, err
	}

	ref, err := git.GlobalGit.LatestTagOrCommit(ctx)
	if err != nil {
		return false, err
	}

	if err := git.GlobalGit.Checkout(ctx, ref); err != nil {
		return false, err
	}
	golog.Info(action)
	return true, nil
}
