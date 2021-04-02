package acceptance

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/cmd/goblog/internal/initialize"
	"github.com/ldelossa/goblog/test"
)

const (
	RelativeSrc = "../../../../goblog"
)

func TestAcceptance(t *testing.T) {
	a := Acceptance{}
	a.BuildWithHijackedEnv(t)
}

// This package provides an acceptance
// test for GoBlog.
//
// An acceptance test acts as a client
// to a GoBlog binary.

// Acceptance is our acceptance test
// and any state held between test
// cases.
type Acceptance struct {
	// the temporary GoBlog
	// home used during the acceptance
	// test.
	tmpHome string
}

// BuildWithHijackedEnv creates a temp GoBlog home directory, build
// a goblog binary in it's ./bin folder, and returns the root of the
// temporary home.
//
// All subsequent invocations of the built goblog binary will utilize
// the $GOBLOG_HOME={a.tmpHome} env variable to point GoBlog to our
// temporary home.
func (a *Acceptance) BuildWithHijackedEnv(t *testing.T) (cleanup func()) {
	cleanup, tmpDir, err := test.HijackEnviroment("bin")
	if err != nil {
		t.Fatalf("failed to create a temporary goblog environment")
	}
	a.tmpHome = tmpDir

	p, err := filepath.Abs(RelativeSrc)
	if err != nil {
		t.Fatalf("failed to resolve abs path for source under test: %v", err)
	}
	err = os.Symlink(p, goblog.Src)
	if err != nil {
		t.Fatalf("failed to link source code under test to goblog.Src: %v", err)
	}

	if _, err := initialize.Build(context.Background()); err != nil {
		t.Fatalf("failed building goblog acceptance binary: %v", err)
	}
	_, err = os.Stat(filepath.Join(tmpDir, "bin", "goblog"))
	if err != nil {
		t.Fatalf("could not stat goblog acceptance binary: %v", err)
	}
	return cleanup
}
