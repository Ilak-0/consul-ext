package git_repo

import (
	"consul-ext/pkg/git_repo/gitea"
	"consul-ext/pkg/git_repo/gitlab"
	"consul-ext/pkg/models"
	gitlabSDK "github.com/xanzy/go-gitlab"
)

type GitFunc interface {
	CreateFile(key string, createFileOptions any) error
	GetFile(key string) (*models.FileData, error)
	UpdateFile(key string, updateFileOptions any) error
	PushCommit(commitOpt any) error
	NewEvent() any
	NewCreateFileOptions(content []byte) any
	NewUpdateFileOptions(content []byte, SHA string) any
	NewCreateCommitOptions() any
	NewCommitActionOptions(content []byte, key string, action gitlabSDK.FileActionValue) any
	GetRepo() models.Repo
}

func APIFunc() GitFunc {
	switch {
	case gitea.Control.Client != nil:
		return gitea.Control
	case gitlab.Control.Client != nil:
		return gitlab.Control
	}
	return nil
}
