package drafts

import (
	"testing"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/test"
)

// TestDraftsListControl is a control
// test ensuring empty local and embedded
// DraftCaches produce empty outputs.
//
// Sanity check.
func TestDraftsListControl(t *testing.T) {
	const (
		want = "ID\tDATE\tTITLE\tSUMMARY\n"
	)
	goblog.LocalPostsCache = []goblog.Post{}
	goblog.EmbeddedPostsCache = []goblog.Post{}

	// test with local first
	got := test.RunWithHijackedStdout(t, list)
	test.CmpEqual(t, got, want)

	// test with embedded
	// we can do this since list just searches
	// for an --embedded arg
	test.HijackOSArgs([]string{"--embedded"})
	got = test.RunWithHijackedStdout(t, list)
	test.CmpEqual(t, got, want)
}

func TestDraftsListLocal(t *testing.T) {
	table := []struct {
		Name string
		N    int
		Want string
	}{
		{Name: "1 Post Listing",
			N:    1,
			Want: "ID\tDATE\t\tTITLE\tSUMMARY\n1\t2021-May-28\t1\t1\n",
		},
		{Name: "5 Post Listing",
			N:    5,
			Want: "ID\tDATE\t\tTITLE\tSUMMARY\n1\t2021-May-28\t5\t5\n2\t2021-May-27\t4\t4\n3\t2021-May-26\t3\t3\n4\t2021-May-25\t2\t2\n5\t2021-May-24\t1\t1\n",
		},
		{Name: "10 Post Listing",
			N:    10,
			Want: "ID\tDATE\t\tTITLE\tSUMMARY\n1\t2021-May-28\t10\t10\n2\t2021-May-27\t9\t9\n3\t2021-May-26\t8\t8\n4\t2021-May-25\t7\t7\n5\t2021-May-24\t6\t6\n6\t2021-May-23\t5\t5\n7\t2021-May-22\t4\t4\n8\t2021-May-21\t3\t3\n9\t2021-May-20\t2\t2\n10\t2021-May-19\t1\t1\n",
		},
	}

	for _, tt := range table {
		t.Run(tt.Name, func(t *testing.T) {
			// hijack local drafts cache
			// so the list commands uses our generated
			// post
			posts := test.GenPosts(tt.N)
			goblog.LocalDraftsCache = posts

			got := test.RunWithHijackedStdout(t, list)
			test.CmpEqual(t, got, tt.Want)
		})
	}
}

func TestDraftsListEmbedded(t *testing.T) {
	table := []struct {
		Name string
		N    int
		Want string
	}{
		{Name: "1 Post Listing",
			N:    1,
			Want: "ID\tDATE\t\tTITLE\tSUMMARY\n1\t2021-May-28\t1\t1\n",
		},
		{Name: "5 Post Listing",
			N:    5,
			Want: "ID\tDATE\t\tTITLE\tSUMMARY\n1\t2021-May-28\t5\t5\n2\t2021-May-27\t4\t4\n3\t2021-May-26\t3\t3\n4\t2021-May-25\t2\t2\n5\t2021-May-24\t1\t1\n",
		},
		{Name: "10 Post Listing",
			N:    10,
			Want: "ID\tDATE\t\tTITLE\tSUMMARY\n1\t2021-May-28\t10\t10\n2\t2021-May-27\t9\t9\n3\t2021-May-26\t8\t8\n4\t2021-May-25\t7\t7\n5\t2021-May-24\t6\t6\n6\t2021-May-23\t5\t5\n7\t2021-May-22\t4\t4\n8\t2021-May-21\t3\t3\n9\t2021-May-20\t2\t2\n10\t2021-May-19\t1\t1\n",
		},
	}

	for _, tt := range table {
		t.Run(tt.Name, func(t *testing.T) {
			// hijack local drafts cache
			// so the list commands uses our generated
			// post
			posts := test.GenPosts(tt.N)
			goblog.EmbeddedDraftsCach = posts
			test.HijackOSArgs([]string{"--embedded"})
			got := test.RunWithHijackedStdout(t, list)
			test.CmpEqual(t, got, tt.Want)
		})
	}
}
