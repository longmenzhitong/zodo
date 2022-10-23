package backup

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"zodo/internal/conf"
	"zodo/internal/files"
)

func Pull() error {
	r, err := git.PlainOpen(files.Dir)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: conf.All.Git.Username,
			Password: conf.All.Git.Password,
		},
	})

	return err
}
