package goblog

import (
	"context"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed drafts/*
var DraftsFS embed.FS

// EmbeddedDraftsCache is a DateSortable slice of posts
// which only contains the metadata of all embedded post.
//
// This allows us to quickly read out date ordered posts
// without walking the embeded filesystem more then once.
//
// initialized in: ./home.go:37
var EmbeddedDraftsCache DateSortable

// LocalDraftsCache is a DateSortable slice of posts
// which only contains the metadata of all local post.
//
// This allows us to quickly read out date ordered posts
// without walking the embeded filesystem more then once.
//
// initialized in: ./home.go:37
var LocalDraftsCache DateSortable

func NewEmbeddedDraftsCache() (DateSortable, error) {
	sorted := DateSortable{}
	err := fs.WalkDir(DraftsFS, "drafts", func(p string, d fs.DirEntry, err error) error {
		if d == nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) != ".post" {
			return nil
		}
		f, err := DraftsFS.Open(p)
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

// walks the "drafts" directory for drafts posts, sorts them by date, and returns
// a list of them.
func NewLocalDraftsCache(ctx context.Context) (DateSortable, error) {
	sorted := DateSortable{}
	err := filepath.Walk(Drafts, func(path string, info os.FileInfo, err error) error {
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

		var post Post

		err = yaml.NewDecoder(f).Decode(&post)
		if err != nil {
			return err
		}
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
