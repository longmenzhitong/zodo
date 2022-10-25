package backup

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"strconv"
	"time"
	"zodo/internal/conf"
	"zodo/internal/files"
	"zodo/internal/todo"
)

const (
	pulledFileName = "pulled"
)

var (
	pulledPath string
)

func init() {
	pulledPath = files.GetPath(pulledFileName)
	files.EnsureExist(pulledPath)
}

func CheckPull() error {
	lines := files.ReadLinesFromPath(pulledPath)
	if len(lines) == 0 {
		return Pull()
	} else if len(lines) == 1 {
		ut, err := strconv.ParseInt(lines[0], 10, 64)
		if err != nil {
			return err
		}
		if time.Unix(ut, 0).Day() < time.Now().Day() {
			return Pull()
		} else {
			return nil
		}
	} else {
		return fmt.Errorf("too many pulled timestamps: %v", lines)
	}
}

func Pull() error {
	// We instantiate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(files.Dir)
	if err != nil {
		return fmt.Errorf("plainopen error: %v", err)
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("worktree error: %v", err)
	}

	// Pull the latest changes from the origin remote and merge into the current branch
	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: conf.All.Git.Username,
			Password: conf.All.Git.Password,
		},
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("pull error: %v", err)
	}

	files.RewriteLinesToPath(pulledPath, []string{strconv.FormatInt(time.Now().Unix(), 10)})

	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil {
		return fmt.Errorf("head error: %v", err)
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return fmt.Errorf("commit object error: %v", err)
	}

	fmt.Println(commit)

	todo.Load()

	return nil
}
