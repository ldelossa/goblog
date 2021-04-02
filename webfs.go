package goblog

import (
	"context"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed web/*
var WebFS embed.FS

func PathSortedLocalWebFiles(ctx context.Context) ([]string, error) {
	sorted := []string{}
	err := filepath.Walk(Web, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		sorted = append(sorted, strings.TrimPrefix(path, Src+"/"))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(sorted)
	return sorted, nil
}

func PathSortedEmbeddedWebFiles(ctx context.Context) ([]string, error) {
	sorted := []string{}
	err := fs.WalkDir(WebFS, "web", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		sorted = append(sorted, strings.TrimPrefix(path, Src+"/"))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(sorted)
	return sorted, nil
}
