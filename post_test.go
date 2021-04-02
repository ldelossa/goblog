package goblog_test

import (
	"sort"
	"testing"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/test"
)

func TestDateSortablePosts(t *testing.T) {
	posts := test.GenRandomPosts(10)
	ds := goblog.DateSortable(posts)
	sort.Sort(ds)

	// ensure each post is always older
	// then the last seen
	{
		var prev goblog.Post
		for i, post := range ds {
			if i == 0 {
				prev = post
				continue
			}
			if prev.Date.Before(post.Date) {
				t.Fatalf("prev: %v, cur: %v", prev.Date, post.Date)
			}
			prev = post
		}
	}

	ps := ds.PathSorted()
	// ensure getting a PathSortableP does not
	// effect our DateSortable array by just
	// doing this again.
	{
		var prev goblog.Post
		for i, post := range ds {
			if i == 0 {
				prev = post
				continue
			}
			if prev.Date.Before(post.Date) {
				t.Fatalf("prev: %v, cur: %v", prev.Date, post.Date)
			}
			prev = post
		}
	}

	// confirm ps is sorted by paths
	{
		var prev *goblog.Post
		for i, post := range ps {
			if i == 0 {
				prev = post
				continue
			}
			if prev.Path < post.Path {
				t.Fatalf("prev: %v, cur: %v", prev.Path, post.Path)
			}
			prev = post
		}
	}
}
