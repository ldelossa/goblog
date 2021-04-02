package ui

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/pkg/golog"
	"gopkg.in/yaml.v3"
)

var ErrNoEditor error = fmt.Errorf("Please set an $EDITOR environment flag")

var ErrNoPath error = fmt.Errorf("The post attempting to be edited has no Path")

// Editor is GoBlog's interface with
// the user's text editor of choice.
//
// Editor wraps a Post and its exported
// methods open the user's editor for
// manipulating it.
//
// Editor utilizes the $EDITOR environment variable
// to resolve which text editor to start.
type Editor struct {
	Post *goblog.Post
	cmd  string
	args []string
}

// NewEditor returns an Editor ready to edit the provided
// post.
func NewEditor(ctx context.Context, post *goblog.Post) (Editor, error) {
	var editor Editor

	env := strings.Split(os.Getenv("EDITOR"), " ")
	if len(env) == 0 {
		return editor, ErrNoEditor
	}

	editor.cmd = env[0]

	// sometimes editor env is more then one
	// command, such as 'code -w' which tells
	// vscode to not background itself.
	if len(env) > 1 {
		editor.args = env[1:]
	}

	editor.Post = post
	return editor, nil
}

func (e Editor) open(ctx context.Context, file string) error {
	// editor arguments + file
	args := make([]string, 0, len(e.args)+1)
	copy(args, e.args)
	args = append(args, file)

	cmd := exec.Command(e.cmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// Edit opens a MarkDown Post in the user's
// editor.
//
// The user actually edits a temporary markdown file
// so editors perform syntax highlighting and linting
// correctly.
//
// Once the user saves the buffer and closes their
// editor the MarkDown contents will be appended
// to the in-memory Post and the date will be updated.
//
// It is expected that the caller persist the
// Post to disk for perminent storage.
//
// It is not done here since the Editor does not
// determine if the draft is published (written to goblog.Posts dir)
// or not (written back to goblog.Drafts dir)
func (e Editor) Edit(ctx context.Context) error {
	// create a temporary markdown file which
	// the user will actually edit.
	if e.Post.Path == "" {
		return ErrNoPath
	}

	mdTmp := strings.ReplaceAll(e.Post.Path, ".post", ".md")
	mdTmp = filepath.Join(goblog.Src, mdTmp)

	f, err := os.OpenFile(mdTmp, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		return err
	}

	_, err = io.WriteString(f, e.Post.MarkDown.Value)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	// open the file in the user's editor
	err = e.open(ctx, mdTmp)
	if err != nil {
		return err
	}

	// read markdown file into a buffer, close fd and remove
	// the draft.
	f, err = os.OpenFile(mdTmp, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	buff, err := io.ReadAll(f)
	if err != nil {
		onErrDumpMarkdown(string(buff), err)
		return err
	}
	err = f.Close()
	if err != nil {
		onErrDumpMarkdown(string(buff), err)
		return err
	}
	err = os.Remove(mdTmp)
	if err != nil {
		onErrDumpMarkdown(string(buff), err)
		return err
	}

	e.Post.MarkDown = yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: yaml.FlowStyle,
		Value: string(buff),
	}
	e.Post.Date = time.Now()
	return nil
}
func onErrDumpMarkdown(md string, err error) {
	golog.Error("Error: failed while saving markdown post: %v", err)
	golog.Error("Dumping your markdown so you don't loose your work...\n")
	fmt.Println("----MARKDOWN BEGIN----")
	fmt.Printf(md)
	fmt.Println("----MARKDOWN END----")
}

// OpenMeta opens the raw Post in the user's editor.
// It's very simple as we are simply opening the file on
// disk for editing and any changes will be saved directly
// to disk.
func (e Editor) EditMeta(ctx context.Context) error {
	path := filepath.Join(goblog.Src, e.Post.Path)
	return e.open(ctx, path)
}
