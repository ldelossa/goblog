package initialize

import (
	"context"
	"log"
	"os"
	"path"
	"testing"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/pkg/dtree"
	"github.com/ldelossa/goblog/test"
)

func TestMain(m *testing.M) {
	cleanup, err := test.HijackEnvironment()
	if err != nil {
		log.Fatalf("TestMain: could not hijack environment")
	}
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func TestInit(t *testing.T) {
	tc := TestCases{}
	tc.Run(t)
}

type TestCases struct {
	// the root of the initialization tree
	root *dtree.Decision
}

func (tc *TestCases) Run(t *testing.T) {
	t.Run("Init", tc.Init)
	t.Run("RemovedHomeDir", tc.RemovedHomeDir)
	t.Run("RemovedSrcDir", tc.RemovedSrcDir)
}

// Init confirms executing the initialization
// decision tree does not fail.
func (tc *TestCases) Init(t *testing.T) {
	tc.root = buildDTree()
	err := tc.root.Execute(context.TODO())
	if err != nil {
		t.Fatalf("failed executing initialization: %v", err)
	}
	checkInitialized(t)
}

// RemovedHomeDir confirms initialization handles the
// deletion of goblog's home directory.
func (tc *TestCases) RemovedHomeDir(t *testing.T) {
	err := os.RemoveAll(goblog.Home)
	if err != nil {
		t.Fatalf("failed to remove home dir for test: %v", err)
	}
	tc.Init(t)
}

// RemovedSrcDir confirms initialization handles the
// deletion of goblog's src code directory.
func (tc *TestCases) RemovedSrcDir(t *testing.T) {
	err := os.RemoveAll(goblog.Src)
	if err != nil {
		t.Fatalf("failed to remove home dir for test: %v", err)
	}
	tc.Init(t)
}

// checkInitialized confirms the GoBlog environment
// is initialized correctly.
func checkInitialized(t *testing.T) {
	_, err := os.Stat(goblog.Src)
	if err != nil {
		t.Fatalf("src directory not found: %v", err)
	}
	_, err = os.Stat(path.Join(goblog.Src, ".git"))
	if err != nil {
		t.Fatalf("src directory is not a git repository: %v", err)
	}
	_, err = os.Stat(path.Join(goblog.Home, "bin"))
	if err != nil {
		t.Fatalf("bin directory not found: %v", err)
	}
	_, err = os.Stat(path.Join(goblog.Home, "bin/goblog"))
	if err != nil {
		t.Fatalf("goblog binary not found: %v", err)
	}
}
