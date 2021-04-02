package drafts

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/test"
)

func TestMain(m *testing.M) {
	cleanup, err := test.HijackEnviroment("drafts")
	if err != nil {
		log.Fatalf("TestMain: could not hijack environment: %v", err)
	}
	code := m.Run()
	cleanup()
	os.Exit(code)
}
func TestDraftsDeleteSuccess(t *testing.T) {
	const (
		want = "Successfully deleted draft 1"
	)
	// gen a single post and create a file
	// at the hijacked Drafts location.
	posts := test.GenPosts(1)
	fqp := filepath.Join(goblog.Src, posts[0].Path)
	f, err := os.Create(fqp)
	if err != nil {
		t.Fatalf("%v", err)
	}
	f.Close()

	// hijack local drafts cache, placing in our generated post
	goblog.LocalDraftsCache = posts

	// hijack args to delete id 1
	test.HijackOSArgs([]string{"goblog", "drafts", "delete", "1"})

	delete(context.Background())

	_, err = os.Stat(fqp)
	if !os.IsNotExist(err) {
		t.Fatalf("file %v still exists", fqp)
	}
}
