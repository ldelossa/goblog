package initialize

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/git"
	"github.com/ldelossa/goblog/pkg/golog"
)

func Build(ctx context.Context) (bool, error) {
	buildDest := path.Join(goblog.Home, "bin")
	action := "Built GoBlog at " + buildDest
	_, err := os.Stat(buildDest)
	switch {
	case os.IsNotExist(err):
		if err := os.Mkdir(buildDest, 0750); err != nil {
			return false, fmt.Errorf("Failed creating bin directory: %w", err)
		}
	case err != nil:
		return false, fmt.Errorf("Failed to stat build directory: %v", err)
	default:
	}
	goPath, err := exec.LookPath("go")
	if err != nil {
		return false, fmt.Errorf("could not find go command in path: %w", err)
	}

	// we will build our GoBlog binary with the tag name or latest
	// commit hash found in GoBlog's src directory.
	//
	// this enables our update detection mechanism, as this tag or commit hash
	// will become the hardcoded constant goblog.Version and compared against
	// the latest upstream tag.
	version, err := git.GlobalGit.HEAD(ctx)
	if err != nil {
		return false, fmt.Errorf("could not GoBlog src tag or commit")
	}
	ldFlag := fmt.Sprintf("-X goblog.Version=%s", version)
	goBuild := exec.Cmd{
		Path:   goPath,
		Args:   []string{"go", "build", "-o", "../bin/goblog", "-ldflags", ldFlag, "./cmd/goblog"},
		Dir:    path.Join(goblog.Home, "src"),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	err = goBuild.Run()
	if err != nil {
		return false, fmt.Errorf("Failed to build GoBlog: %v", err)
	}
	golog.Info(action)
	return true, nil
}
