package goblog

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/google/go-cmp/cmp"
)

// Differ creates diff reports for our Posts,
// Drafts, Web, and Config directories/files.
//
// Differ is used in coordination with go-cmp
// and fulfills an interface.
//
// The only methods interesting for GoBlog's use
// cases are the ones generating diffs for our
// directories and files.
type Differ struct {
	path     cmp.Path
	diffs    []string
	withDiff bool
	fs       embed.FS
}

// reset sets the differ back to a zero state
func (r *Differ) reset() {
	r.path = nil
	r.diffs = nil
	r.withDiff = false
	r.fs = embed.FS{}
}

// Drafts will return a utf-8 encoded diff string indicating
// any diffs between local and embedded drafts.
//
// if withDiff is true and two drafts with the same name differ in content
// the differing lines will be displayed.
//
// An empty string indicates both the local and embedded contents are the
// the same.
func (r *Differ) Drafts(ctx context.Context, withDiff bool) (string, error) {
	r.reset()

	localPaths, embeddedPaths := []string{}, []string{}

	local := LocalDraftsCache.PathSorted()
	for _, post := range local {
		localPaths = append(localPaths, post.Path)
	}
	embedded := EmbeddedDraftsCache.PathSorted()
	for _, post := range embedded {
		embeddedPaths = append(embeddedPaths, post.Path)
	}

	r.fs = DraftsFS
	r.withDiff = withDiff

	cmp.Equal(embeddedPaths, localPaths, cmp.Reporter(r))

	diff := r.String()

	return diff, nil
}

// Posts will return a utf-8 encoded diff string indicating
// any diffs between local and embedded posts.
//
// if withDiff is true and two posts with the same name differ in content
// the differing lines will be displayed.
//
// An empty string indicates both the local and embedded contents are the
// the same.
func (r *Differ) Posts(ctx context.Context, withDiff bool) (string, error) {
	r.reset()

	localPaths, embeddedPaths := []string{}, []string{}

	local := LocalPostsCache.PathSorted()
	for _, post := range local {
		localPaths = append(localPaths, post.Path)
	}

	embedded := EmbeddedPostsCache.PathSorted()
	for _, post := range embedded {
		embeddedPaths = append(embeddedPaths, post.Path)
	}

	r.fs = PostsFS
	r.withDiff = withDiff

	cmp.Equal(embeddedPaths, localPaths, cmp.Reporter(r))

	diff := r.String()

	return diff, nil
}

// Config will return a utf-8 encoded diff string indicating
// any diffs between the local and embedded config file.
//
// If withDiff is true and two files with the same name differ in content
// the differing lines will be displayed.
//
// An empty string indicates both the local and embedded contents are the
// the same.
func (r *Differ) Config(ctx context.Context) (string, error) {
	r.reset()

	var localBuff bytes.Buffer
	var embeddedBuff bytes.Buffer

	local, err := os.Open(
		filepath.Join(Configs, "config.yaml"),
	)
	if err != nil {
		return "", fmt.Errorf("failed opening local config: %v", err)
	}

	_, err = io.Copy(&localBuff, local)
	if err != nil {
		return "", fmt.Errorf("failed copying local config: %v", err)
	}
	local.Close()

	embed, err := ConfigFS.Open("config/config.yaml")
	if err != nil {
		panic("could not open config: " + err.Error())
	}

	_, err = io.Copy(&embeddedBuff, embed)
	if err != nil {
		return "", fmt.Errorf("failed copying embedded config: %v", err)
	}
	embed.Close()

	return cmp.Diff(embeddedBuff.String(), localBuff.String()), nil
}

// Web will return a utf-8 encoded diff string indicating any
// diffs between the local and embedded web root content.
//
// An empty string indicates both the local and embedded contents are the
// the same.
func (r *Differ) Web(ctx context.Context, withDiff bool) (string, error) {
	r.reset()

	var err error

	localPaths, embeddedPaths := []string{}, []string{}

	localPaths, err = PathSortedLocalWebFiles(ctx)
	if err != nil {
		return "", err
	}

	embeddedPaths, err = PathSortedEmbeddedWebFiles(ctx)
	if err != nil {
		return "", err
	}

	r.fs = WebFS
	r.withDiff = withDiff

	cmp.Equal(embeddedPaths, localPaths, cmp.Reporter(r))

	diff := r.String()

	return diff, nil
}

func (r *Differ) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *Differ) String() string {
	return strings.Join(r.diffs, "\n")
}

func (r *Differ) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *Differ) Report(rs cmp.Result) {
	// if file lives both in embedded and local tree
	// diff the files themselves to determine equality.
	if rs.Equal() {
		vx, _ := r.path.Last().Values()
		// check if local and embedded files are the same.
		file := fmt.Sprintf("%v", vx)

		embedded, err := r.fs.Open(file)
		if err != nil {
			color.Red(`
Error: failed to open embedded file for diffing: %v

`, err)
			return
		}
		local, err := os.OpenFile(
			filepath.Join(Src, file),
			os.O_RDONLY,
			0,
		)
		if err != nil {
			color.Red(`
Error: failed to open local file for diffing: %v

`, err)
			return
		}

		embedBuff, err := io.ReadAll(embedded)
		if err != nil {
			color.Red(`
Error: failed to read embedded file for diffing into buffer: %v

`, err)
			return
		}
		localBuff, err := io.ReadAll(local)
		if err != nil {
			color.Red(`
Error: failed to read local file for diffing into buffer: %v

`, err)
			return
		}

		diff := cmp.Diff(string(embedBuff), string(localBuff))
		if diff != "" {
			if r.withDiff {
				r.diffs = append(r.diffs, fmt.Sprintf("\t+ %+v dirty\n%s\n", vx, diff))
			} else {
				r.diffs = append(r.diffs, fmt.Sprintf("\t+ %+v dirty\n", vx))
			}
		}
	}

	if !rs.Equal() {
		vx, vy := r.path.Last().Values()
		switch {
		case vx.IsValid() && vy.IsValid():
			r.diffs = append(r.diffs, fmt.Sprintf("\t- %+v\n+ %+v\n", vx, vy))
		case vx.IsValid() && !vy.IsValid():
			r.diffs = append(r.diffs, fmt.Sprintf("\t- %+v\n", vx))
		case !vx.IsValid() && vy.IsValid():
			r.diffs = append(r.diffs, fmt.Sprintf("\t+ %+v\n", vy))
		default:
			return
		}
	}
}
