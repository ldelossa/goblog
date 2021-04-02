package test

import (
	"math/rand"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ldelossa/goblog"
)

// RandomString returns a random
// string utilizing a-z of length
// n.
func RandomString(n int) string {
	chars := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		chars = append(chars, byte(0x61+(rand.Int()%26)))
	}
	return string(chars)
}

func GenRandomPosts(n int) []goblog.Post {
	posts := make([]goblog.Post, 0, n)
	for i := 0; i < n; i++ {
		durationMod := rand.Int() % 121
		duration := time.Duration(-durationMod) * (24 * time.Hour)
		posts = append(posts, goblog.Post{
			Path:    filepath.Join(goblog.Posts, RandomString(4)) + ".post",
			Hero:    filepath.Join(goblog.Posts, RandomString(4)) + ".png",
			Title:   RandomString(4) + " " + RandomString(4),
			Summary: RandomString(4) + " " + RandomString(4),
			Date:    time.Now().Add(duration),
		})
	}
	return posts
}

// GenPosts will generate predictable posts for
// testing.
//
// Posts in GoBlog are currently sorted by either
// Path or Date.
//
// Posts are returned with both their Date (newest->oldest)
// and Paths (lexical highest->lowest) in descending order.
//
// This is useful to confirm sorting algorithms work, as sorting
// by paths should reverse the list, and sorting again by date
// shoud return the list its previous state.
func GenPosts(n int) []goblog.Post {
	posts := make([]goblog.Post, 0, n)
	now := time.Now()
	for i := 0; i < n; i++ {
		s := strconv.Itoa(n - i)
		posts = append(posts, goblog.Post{
			Path:    filepath.Join("drafts", s),
			Title:   s,
			Summary: s,
			Date:    now.Add(-time.Duration(i) * (24 * time.Hour)),
		})
	}
	return posts
}
