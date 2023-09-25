package models

import "code.gitea.io/sdk/gitea"

const (
	GITLAB = Repo("gitlab")
	GITEA  = Repo("gitea")
)

type Repo string

type FileData struct {
	Content string
	SHA     string
}

type GiteaPayloadCommit struct {
	Commits []*gitea.PayloadCommit `json:"commits"`
}
