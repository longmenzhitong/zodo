package backup

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
	"time"
	"zodo/internal/conf"
	"zodo/internal/cst"
	"zodo/internal/files"
)

func CheckPull() error {
	todayPulled := getPulledPath(time.Now())
	if _, err := os.Stat(todayPulled); err == nil {
		return nil
	}
	err := Pull()
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}
	files.EnsureExist(todayPulled)

	yesterdayPulled := getPulledPath(time.Now().AddDate(0, 0, -1))
	if _, err = os.Stat(yesterdayPulled); err == nil {
		err = os.Remove(yesterdayPulled)
		if err != nil {
			return err
		}
	}
	return nil
}

func getPulledPath(t time.Time) string {
	return files.GetPath(fmt.Sprintf("%s.pulled", t.Format(cst.LayoutMonthDay)))
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
	if err != nil {
		return fmt.Errorf("pull error: %v", err)
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
