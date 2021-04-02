package ui

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/ldelossa/goblog"
)

// GlobalPrompter is a default
// Prompter attached to Stdin and
// Stdout
var GlobalPrompter = NewPrompter()

// ErrFlush indicates an error while flushing
// a Prompter's buffer to stdout.
type ErrFlush struct {
	err error
}

func (e ErrFlush) Error() string {
	return e.err.Error()
}

// ErrRead indicates an error while reading
// a Prompter's stdin buffer.
type ErrRead struct {
	err error
}

func (e ErrRead) Error() string {
	return e.err.Error()
}

// Prompter exports methods for interacting
// with the user via stdin and stdout prompts.
//
// Prompter is a bufio.ReadWriter attached
// to the standard pipes.
type Prompter struct {
	// a ReadWriter which reads from Stdint
	// and writes to Stdou
	stdio bufio.ReadWriter
}

func NewPrompter() Prompter {
	rw := bufio.ReadWriter{
		Reader: bufio.NewReader(os.Stdin),
		Writer: bufio.NewWriter(os.Stdout),
	}
	return Prompter{
        rw,
    }
}

// DraftBuilder issues several posts to the user to construct
// the Metadata of a draft Post.
func (p Prompter) DraftBuilder(ctx context.Context) (goblog.Post, error) {
	const (
		TitlePrompt   = "What's the title of this post (required)?\n> "
		SummaryPrompt = "What's the summary of this post(required)?\n> "
		HeroPrompt    = "Path to a hero image.\nHero images live in the /posts directory so provide a path such as '/posts/myposthero.png'\nType 'none' for no hero image.\n> "
	)

	var post goblog.Post
	var err error

	post.Title, err = p.prompt(ctx, TitlePrompt)
	if err != nil {
		return post, err
	}

	post.Summary, err = p.prompt(ctx, SummaryPrompt)
	if err != nil {
		return post, err
	}

	post.Hero, err = p.prompt(ctx, HeroPrompt)
	if err != nil {
		return post, err
	}

	formatted := strings.ReplaceAll(post.Title, " ", "_")
	formatted = strings.ToLower(formatted)
	// this is a temporary path used to aide building a
	// draft.
	//
	// you'll notice real paths are set when we initialize
	// our in-memory caches.
	//
	// see: goblog/postsfs.go:66 as an example.
	post.Path = formatted + ".post"
	return post, nil
}

// ShouldPublishPost asks the user if they want to publish a
// draft post and returns a bool indicating their decision.
func (p Prompter) ShouldPublishPost(ctx context.Context) (bool, error) {
	const (
		prompt = "Publish this post? ('true', 'false')\n> "
	)
	return p.boolPrompt(ctx, prompt)
}

// ShouldBuild asks the user if they want to build a new goblog
// binary and returns a bool indicating their decision.
func (p Prompter) ShouldBuild(ctx context.Context) (bool, error) {
	const (
		prompt = "Build a new GoBlog binary? ('true', 'false')\n> "
	)
	return p.boolPrompt(ctx, prompt)
}

// prompt will prompt the user and block until a non-empty
// string response is read.
func (p Prompter) prompt(ctx context.Context, prompt string) (string, error) {
	var s string
	for s == "" {
		p.stdio.WriteString(prompt)
		err := p.stdio.Flush()
		if err != nil {
			return "", ErrFlush{err}
		}
		s, err = p.stdio.Reader.ReadString('\n')
		s = strings.TrimSpace(s)
	}
	return s, nil
}

// boolPrompt will prompt the user and block until a non-empty
// boolean response is read.
func (p Prompter) boolPrompt(ctx context.Context, prompt string) (bool, error) {
	var b *bool = nil
	for b == nil {
		s, err := p.prompt(ctx, prompt)
		bb, err := strconv.ParseBool(s)
		if err != nil {
			continue
		}
		b = &bb
	}
	return *b, nil
}
