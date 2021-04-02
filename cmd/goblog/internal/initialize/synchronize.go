package initialize

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/pkg/golog"
)

// Synchronize synchronizes the current
// binary's embedded content with the local src
// tree.
//
// This decision always calls its Yes branch or
// errors out.
func Synchronize(ctx context.Context) (bool, error) {
	action := "Synchronized local and embedded content"
	differ := goblog.Differ{}
	shouldLog := false

	// sync drafts
	draftsDiff, err := differ.Drafts(ctx, false)
	if err != nil {
		return false, err
	}

	if draftsDiff != "" {
		shouldLog = true
		for _, post := range goblog.EmbeddedDraftsCache {
			path := filepath.Join(goblog.Src, post.Path)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				embedded, err := goblog.DraftsFS.Open(post.Path)
				if err != nil {
					return false, err
				}

				f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0660)
				if err != nil {
					return false, err
				}

				_, err = io.Copy(f, embedded)
				if err != nil {
					return false, err
				}

				f.Close()
				embedded.Close()
			}
		}
	}

	// sync posts
	postsDiff, err := differ.Posts(ctx, false)
	if err != nil {
		return false, err
	}

	if postsDiff != "" {
		shouldLog = true
		for _, post := range goblog.EmbeddedPostsCache {
			path := filepath.Join(goblog.Src, post.Path)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				embedded, err := goblog.PostsFS.Open(post.Path)
				if err != nil {
					return false, err
				}

				f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0660)
				if err != nil {
					return false, err
				}

				_, err = io.Copy(f, embedded)
				if err != nil {
					return false, err
				}

				f.Close()
				embedded.Close()
			}
		}
	}

	// sync config
	confDiff, err := differ.Config(ctx)
	if err != nil {
		return false, err
	}

	if confDiff != "" {
		shouldLog = true
		if _, err := os.Stat(filepath.Join(goblog.Configs, "config.yaml")); os.IsNotExist(err) {
			src, err := goblog.ConfigFS.Open("config/config.yaml")
			if err != nil {
				return false, fmt.Errorf("failed opening emedded config file: %w", err)
			}
			dest, err := os.OpenFile(
				filepath.Join(goblog.Configs, "config.yaml"),
				os.O_CREATE|os.O_WRONLY,
				0660,
			)
			if err != nil {
				return false, fmt.Errorf("failed opening local config file: %w", err)
			}
			_, err = io.Copy(dest, src)
			if err != nil {
				return false, fmt.Errorf("failed writing local config file: %w", err)
			}
			src.Close()
			dest.Close()
		}
	}

	// sync web dir
	webDiff, err := differ.Web(ctx, false)
	if err != nil {
		return false, err
	}

	if webDiff != "" {
		shouldLog = true
		err = fs.WalkDir(goblog.WebFS, "web", func(p string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				if _, err := os.Stat(filepath.Join(goblog.Src, p)); os.IsNotExist(err) {
					err := os.MkdirAll(p, 0770)
					if err != nil && !os.IsExist(err) {
						return fmt.Errorf("failed creating dir in web root: %w", err)
					}
				}
			}

			if _, err := os.Stat(filepath.Join(goblog.Src, p)); os.IsNotExist(err) {
				src, err := goblog.WebFS.Open(p)
				if err != nil {
					return fmt.Errorf("failed opening emedded web file %s: %w", p, err)
				}
				defer src.Close()
				dest, err := os.OpenFile(
					filepath.Join(goblog.Src, p),
					os.O_CREATE|os.O_WRONLY,
					0660,
				)
				if err != nil {
					return fmt.Errorf("failed opening local web file %s: %w", p, err)
				}
				defer dest.Close()
				_, err = io.Copy(dest, src)
				if err != nil {
					return fmt.Errorf("failed writing local web file %s: %w", p, err)
				}
			}
			return nil
		})
	}
	if shouldLog {
		golog.Info(action)
	}
	return true, nil
}
