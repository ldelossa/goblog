package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ldelossa/goblog"
)

var GlobalGit = NewGit()

type ErrFetch struct {
	err error
}

func (e ErrFetch) Error() string {
	return fmt.Sprintf("Error while git Fetch: %v", e.err.Error())
}

type ErrClone struct {
	err error
}

func (e ErrClone) Error() string {
	return fmt.Sprintf("Error while git Clone: %v", e.err.Error())
}

type ErrReset struct {
	err error
}

func (e ErrReset) Error() string {
	return fmt.Sprintf("Error while git Reset: %v", e.err.Error())
}

type ErrClean struct {
	err error
}

func (e ErrClean) Error() string {
	return fmt.Sprintf("Error while git Clean: %v", e.err.Error())
}

type ErrCheckout struct {
	err error
}

func (e ErrCheckout) Error() string {
	return fmt.Sprintf("Error while git Checkout: %v", e.err.Error())
}

type ErrLatestTagOrCommit struct {
	err error
}

func (e ErrLatestTagOrCommit) Error() string {
	return fmt.Sprintf("Error while git LatestTagOrCommit: %v", e.err.Error())
}

type ErrHEAD struct {
	err error
}

func (e ErrHEAD) Error() string {
	return fmt.Sprintf("Error while git HEAD: %v", e.err.Error())
}

// Git is a opinionated synchronous
// git runner for use with the
// goblog.Src directory.
type Git struct {
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func NewGit() Git {
	return Git{
		stdout: &bytes.Buffer{},
		stderr: &bytes.Buffer{},
	}
}

// reset resets the internal buffers
// capturing stdin and stdout of the last
// git command ran.
func (g Git) reset() {
	g.stdout.Reset()
	g.stderr.Reset()
}

// run is a helper which runs a git process,
// captures its stdout and stderr, and
// resets the buffer once std out is
// communicated.
//
// the run helper performs all git actions
// against goblog.Src
//
// any stdout is returned with all leading and
// trailing white space removed.
func (g Git) run(args []string) (string, error) {
	defer g.reset()
	cmd := exec.Command("git", args...)
	cmd.Stdout = g.stdout
	cmd.Stderr = g.stderr
	cmd.Dir = goblog.Src
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Err: %v Stderr: %v", err, g.stderr.String())
	}
	return strings.TrimSpace(g.stdout.String()), nil
}

// Fetch fetches new tags from GoBlog's upstream
// repo.
func (g Git) Fetch(ctx context.Context) error {
	args := []string{"fetch", "--tags"}
	if _, err := g.run(args); err != nil {
		return ErrFetch{err}
	}
	return nil
}

// Clone clones the GoBlog remote to the configured
// goblog.Src directory.
func (g Git) Clone(ctx context.Context, remote string) error {
	// clone does not use the run helper, since
	// the run helper assumes goblog.Src already
	// exists and sets the working dir to goblog.Home.
	defer g.reset()
	args := []string{"clone", remote, "src"}
	cmd := exec.Command("git", args...)
	cmd.Stdout = g.stdout
	cmd.Stderr = g.stderr
	cmd.Dir = goblog.Home
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Err: %v Stderr: %v", err, g.stderr.String())
	}
	return nil
}

// Reset performs a git reset to the upstream
// origin, effectively moving goblog.Src to the
// latest commit.
func (g Git) Reset(ctx context.Context) error {
	args := []string{"reset", "--hard", "origin/master"}
	if _, err := g.run(args); err != nil {
		return ErrReset{err}
	}
	return nil
}

// Clean will remove any files not in the upstream
// repository
func (g Git) Clean(ctx context.Context) error {
	args := []string{"clean", "-xdf"}
	if _, err := g.run(args); err != nil {
		return ErrClean{err}
	}
	return nil
}

// Checkout will check goblog.Src to a particular
// commit-ish
func (g Git) Checkout(ctx context.Context, commit string) error {
	// fetch first
	if err := g.Fetch(ctx); err != nil {
		return err
	}
	args := []string{"checkout", commit}
	if _, err := g.run(args); err != nil {
		return ErrCheckout{err}
	}
	return nil
}

// HEAD returns the short commit hash of goblog.Src's
// current head ref.
func (g Git) HEAD(ctx context.Context) (string, error) {
	// first try to return a tag if there's one.
	args := []string{"describe", "HEAD", "--tag", "--candidates=0"}
	out, err := g.run(args)
	if err != nil {
		return out, nil
	}

	// a tag didn't exist, return the commit hash
	args = []string{"rev-parse", "--short", "HEAD"}
	out, err = g.run(args)
	if err != nil {
		return "", ErrHEAD{err}
	}
	return out, nil
}

// LatestTagOrCommit will return the latest tag or the latest
// commit if no tag exists.
func (g Git) LatestTagOrCommit(ctx context.Context) (string, error) {
	// fetch first
	if err := g.Fetch(ctx); err != nil {
		return "", err
	}
	args := []string{"describe", "--tags", "--abbrev=0", "--always"}
	out, err := g.run(args)
	if err != nil {
		return "", ErrLatestTagOrCommit{err}
	}
	return out, nil
}

// Same as LatestTagOrCommit but simply returns an empty
// string if no tag can be found.
func (g Git) LatestTag(ctx context.Context) (string, error) {
	tag, err := g.LatestTagOrCommit(ctx)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(tag, "v") {
		return "", nil
	}
	return tag, nil
}
