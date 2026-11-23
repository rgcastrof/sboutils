package pkg

import (
	"errors"
	"io/fs"
	"os"

	"github.com/rgcastrof/sboutils/cli"
	"github.com/rgcastrof/sboutils/internal/config"
)

type Repo struct{}

func (r *Repo) pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (r *Repo) Sync(cfg *config.Config) error {
	ex := cli.NewExecutor()
	var args []string
	if exists, _ := r.pathExists(cfg.RepoPath); exists {
		args = cfg.PullArgs()
	} else {
		args = cfg.SyncArgs()
	}

	return ex.Run(cfg.GitBin, args)
}

func (r *Repo) Clean(cfg *config.Config) error {
	ex := cli.NewExecutor()
	return ex.Run(cfg.GitBin, cfg.CleanArgs())
}
