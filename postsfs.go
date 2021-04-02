package goblog

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed posts/*
var PostsFS embed.FS

// EmbeddedPostsCache is a DateSortable slice of posts
// which only contains the metadata of all embedded post.
//
// This allows us to quickly read out date ordered posts
// without walking the embeded filesystem more then once.
//
// initialized in: ./home.go:37
var EmbeddedPostsCache DateSortable

// LocalPostsCache is a DateSortable slice of posts
// which only contains the metadata of all local post.
//
// This allows us to quickly read out date ordered posts
// without walking the local filesystem more then once.
//
// initialized in: ./home.go:37
var LocalPostsCache DateSortable

func NewEmbeddedPostsCache() (DateSortable, error) {
	sorted := DateSortable{}
	err := fs.WalkDir(PostsFS, "posts", func(p string, d fs.DirEntry, err error) error {
		if d == nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) != ".post" {
			return nil
		}
		f, err := PostsFS.Open(p)
		if err != nil {
			return err
		}

		var post Post
		err = yaml.NewDecoder(f).Decode(&post)
		if err != nil {
			return err
		}
		// .empty is a trick to embed an "empty" posts directory
		// as embed fs api require at least a file inside a dir
		// its embedding. we will just ignore it.
		if post.Title == "_empty" {
			return nil
		}

		sorted = append(sorted, Post{
			Path:    p,
			Title:   post.Title,
			Summary: post.Summary,
			Date:    post.Date,
			Hero:    post.Hero,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(sorted) > 0 {
		sort.Sort(sorted)
	}
	return sorted, nil
}

// walks the local "posts" directory for local posts, sorts them by date, and returns
// a list of them.
func NewLocalPostsCache(ctx context.Context) (DateSortable, error) {
	sorted := DateSortable{}
	err := filepath.Walk(Posts, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		if filepath.Ext(f.Name()) != ".post" {
			return nil
		}

		var post Post
		err = yaml.NewDecoder(f).Decode(&post)
		if err != nil {
			return fmt.Errorf("failed reading post %v: %v", f.Name(), err)
		}
		// ignore the empty post
		if post.Title == "_empty" {
			return nil
		}
		post.Path = strings.TrimPrefix(path, Src+"/")
		sorted = append(sorted, post)
		return nil
	})
	if err != nil {
		return sorted, err
	}
	sort.Sort(sorted)
	return sorted, nil
}
