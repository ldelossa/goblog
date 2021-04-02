package test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"
	"testing"

	"github.com/ldelossa/goblog"
)

// HijackEnviroment re-writes the global goblog variables
// redirecting its environment to a temporary
// folder.
//
// A cleanup function is returned to remove the tmp
// dir.
func HijackEnviroment(mkdir ...string) (func(), string, error) {
	// hijack goblog.Home for isolated testing.
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, "", fmt.Errorf("test main: failed to create tmp home: %v", err)
	}
	goblog.Home = tmpDir
	goblog.Src = path.Join(goblog.Home, "src")
	goblog.Posts = path.Join(goblog.Src, "posts")
	goblog.Drafts = path.Join(goblog.Src, "drafts")
	goblog.Web = path.Join(goblog.Src, "web")
	goblog.Configs = path.Join(goblog.Src, "config")

	for _, dir := range mkdir {
		var err error
		switch dir {
		case "bin":
			err = os.MkdirAll(filepath.Join(goblog.Home, "bin"), 0777)
		case "src":
			err = os.MkdirAll(goblog.Src, 0777)
		case "posts":
			err = os.MkdirAll(goblog.Posts, 0777)
		case "drafts":
			err = os.MkdirAll(goblog.Drafts, 0777)
		case "web":
			err = os.MkdirAll(goblog.Web, 0777)
		case "config":
			err = os.MkdirAll(goblog.Configs, 0777)
		default:
		}
		if err != nil {
			os.RemoveAll(tmpDir)
			return nil, "", err
		}
	}

	return func() { os.RemoveAll(tmpDir) }, tmpDir, nil
}

func HijackOSArgs(args []string) {
	os.Args = args
}

type LockableBuffer struct {
	sync.Mutex
	bytes.Buffer
}

func RunWithHijackedStdout(t *testing.T, cmd func(context.Context)) string {
	var buf LockableBuffer
	ctx, cancel := context.WithCancel(context.Background())

	err := HijackStdout(ctx, &buf)
	if err != nil {
		t.Fatalf("%v", err)
	}

	cmd(ctx)
	cancel()

	buf.Lock()
	defer buf.Unlock()
	return buf.String()
}

// HijackStdout copies from stdout into the supplied lockable
// buffer.
func HijackStdout(ctx context.Context, buf *LockableBuffer) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	old := os.Stdout
	os.Stdout = w

	// locked helps HijackStdout block
	// returning to the client until
	// the lockable buffer is indeed locked.
	locked := make(chan struct{})
	// Lock the buffer, any access
	// other access to this buffer
	// will block on the lock access.
	go func() {
		buf.Lock()
		// tell HijackStdout buffer is locked,
		// cool to return to caller
		locked <- struct{}{}

		// copy data until w.Close()
		io.Copy(buf, r)

		// swap stdout back
		os.Stdout = old
		buf.Unlock()
	}()

	// Use ctx.Done() as an indication
	// that testing code is finished
	// writing to stdout.
	//
	// When this occurs stdout is swapped back,
	// W is closed which will cause a Pipe error
	// in the above GoRoutine, and subsequentlly
	// the buffer will be unlocked.
	go func() {
		select {
		case <-ctx.Done():
			w.Close()
		}
	}()
	<-locked
	return nil
}
