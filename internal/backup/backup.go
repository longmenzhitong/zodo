package backup

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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

func Push() error {
	// Opens an already existing repository.
	r, err := git.PlainOpen(files.Dir)
	if err != nil {
		return fmt.Errorf("plainopen error: %v", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("worktree error: %v", err)
	}

	// Adds the new file to the staging area.
	_, err = w.Add(".")
	if err != nil {
		return fmt.Errorf("add error: %v", err)
	}

	// We can verify the current status of the worktree using the method Status.
	status, err := w.Status()
	if err != nil {
		return fmt.Errorf("status error: %v", err)
	}

	fmt.Println(status)

	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit Since version 5.0.1, we can omit the Author signature, being read
	// from the git config files.
	commit, err := w.Commit("-", &git.CommitOptions{
		Author: &object.Signature{
			Name:  conf.All.Git.Username,
			Email: conf.All.Git.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("commit error: %v", err)
	}

	// Prints the current HEAD to verify that all worked well.
	obj, err := r.CommitObject(commit)
	if err != nil {
		return fmt.Errorf("commit object error: %v", err)
	}

	fmt.Println(obj)

	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: conf.All.Git.Username,
			Password: conf.All.Git.Password,
		},
	})
	if err != nil {
		return fmt.Errorf("push error: %v", err)
	}

	return nil
}
