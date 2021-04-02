package initialize

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/pkg/golog"
)

// CheckGitRepo is a Decision which
// determines if GoBlog's source code exist in GoBlog's
// home directory.
//
// If it does it calls its Yes branch, if not it
// calls its No branch.
func CheckGitRepo(ctx context.Context) (bool, error) {
	action := "Can't find GoBlog's source code"
	src := path.Join(goblog.Home, "src", ".git")
	info, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			golog.Info(action)
			return false, nil
		}
		return false, fmt.Errorf("Failed to check if GoBlog source code is available: %w", err)
	}
	if !info.IsDir() {
		golog.Info(action)
		return false, nil
	}
	return true, nil
}
