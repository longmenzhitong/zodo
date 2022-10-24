package backup

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
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
		return err
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		return err
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
		return err
	}

	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil {
		return err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	fmt.Println(commit)

	return nil
}
