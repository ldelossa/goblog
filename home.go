package goblog

import (
	"context"
	"os"
	"os/user"
	"path"

	"github.com/fatih/color"
)

var (
	// Home is the root folder GoBlog works out of.
	// It provides a well defined home enviroment
	// where other directories can root themselves.
	//
	// Resolving a desired Home directory is dependent
	// on user.Current() returning the current user.
	Home string
	// Src is a directory nested in Home where GoBlog's
	// downstream (forked) source code lives.
	Src string
	// Posts is a directory which hold published
	// GoBlog posts
	Posts string
	// Drafts is a directory which holds draft blog
	// posts until they are published.
	Drafts string
	// Web is a directory which holds a user's front
	// end web application.
	Web string
	// Configs is a directory which holds GoBlog's embedded
	// configuration.
	Configs string
)

func init() {
	// For testing purposes, if we find the demonstrated
	// environment variable, use it as our GoBlog home.
	// see: cmd/goblog/acceptance
	var homeDir string
	if env := os.Getenv("GOBLOG_HOME"); env != "" {
		homeDir = env
	} else {
		// we need to be able to determine the current user
		// to resolve home dirs.
		usr, err := user.Current()
		if err != nil {
			color.Red("Error: GoBlog must be able to determine the current user.")
			os.Exit(1)
		}
		homeDir = usr.HomeDir
	}
	Home = path.Join(homeDir, "goblog")
	Src = path.Join(Home, "src")
	Posts = path.Join(Src, "posts")
	Drafts = path.Join(Src, "drafts")
	Web = path.Join(Src, "web")
	Configs = path.Join(Src, "config")

	var err error
	EmbeddedPostsCache, err = NewEmbeddedPostsCache()
	if err != nil {
		panic("could not create PostsCache: " + err.Error())
	}
	LocalPostsCache, err = NewLocalPostsCache(context.Background())
	if err != nil {
		panic("could not create local PostsCache: " + err.Error())
	}
	EmbeddedDraftsCache, err = NewEmbeddedDraftsCache()
	if err != nil {
		panic("could not create DSCache: " + err.Error())
	}
	LocalDraftsCache, err = NewLocalDraftsCache(context.Background())
	if err != nil {
		panic("could not create DSCache: " + err.Error())
	}
}
