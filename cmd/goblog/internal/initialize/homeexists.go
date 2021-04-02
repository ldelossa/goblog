package initialize

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ldelossa/goblog"
	"github.com/ldelossa/goblog/pkg/golog"
)

// HomeExists determines if GoBlog's home directory
// exists. This home directory is resolved via the user.Current() function
// in the os/user package.
//
// If its home does exist the decision calls its Yes branch and logs its location.
//
// If it does not exist an attempt to create the dir is made.
// On success its No branch will be called indicating a home directory
// did not exist before this Decision.
//
// If an error occurs creating the directory an error value is returned.
func HomeExists(ctx context.Context) (bool, error) {
	action := "Created GoBlog's home directory at " + goblog.Home
	fi, err := os.Stat(goblog.Home)
	pathErr := new(os.PathError)
	switch {
	case err == nil:
		if !fi.IsDir() {
			return false, fmt.Errorf("Looks like you have a regular file named goblog in your home dir: %v. You'll need to remove this before GoBlog can continue.", goblog.Home)
		}
		return true, nil
	case errors.As(err, &pathErr):
	default:
		return false, err
	}

	err = os.Mkdir(goblog.Home, 0o750)
	if err != nil {
		return false, fmt.Errorf("Error creating GoBlog home: %w", err)
	}
	golog.Info(action)
	return false, nil
}
