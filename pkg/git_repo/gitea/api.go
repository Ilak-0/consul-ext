package gitea

import (
	"code.gitea.io/sdk/gitea"
	"consul-ext/pkg/config"
	"consul-ext/pkg/models"
	"encoding/base64"
	"fmt"
	gitlabSDK "github.com/xanzy/go-gitlab"
	"time"
)

func (g *Opt) UpdateFile(key string, updateFileOptions any) error {
	g.Lock.Lock()
	defer g.Lock.Unlock()
	options := updateFileOptions.(gitea.UpdateFileOptions)
	_, _, err := g.Client.UpdateFile(config.Data.Gitea.Owner, config.Data.Gitea.Repo, key, options)
	if err != nil {
		fmt.Println("Error updating file content in Gitea:", key, err)
		return err
	}
	return nil
}

func (g *Opt) GetFile(key string) (*models.FileData, error) {
	g.Lock.RLock()
	defer g.Lock.RUnlock()
	data, _, err := g.Client.GetContents(config.Data.Gitea.Owner, config.Data.Gitea.Repo, config.Data.Gitea.Branch, key)
	if err != nil {
		return nil, err
	}
	return &models.FileData{
		Content: *data.Content,
		SHA:     data.SHA,
	}, nil
}

func (g *Opt) CreateFile(key string, createFileOptions any) error {
	g.Lock.Lock()
	defer g.Lock.Unlock()
	options := createFileOptions.(gitea.CreateFileOptions)
	_, _, err := g.Client.CreateFile(config.Data.Gitea.Owner, config.Data.Gitea.Repo, key, options)
	if err != nil {
		fmt.Println("Error creating file in Gitea:", key, options.Content, err)
		return err
	}
	return nil
}

func (g *Opt) NewEvent() any {
	return new(models.GiteaPayloadCommit)
}

func (g *Opt) NewCreateFileOptions(content []byte) any {
	return gitea.CreateFileOptions{
		FileOptions: gitea.FileOptions{
			Message:    time.Now().String(),
			BranchName: config.Data.Gitea.Branch,
			Author: gitea.Identity{
				Name: "consul-control",
			},
			Committer: gitea.Identity{
				Name: "consul-control",
			},
			Signoff: true,
		},
		Content: base64.StdEncoding.EncodeToString(content),
	}
}

func (g *Opt) NewUpdateFileOptions(content []byte, SHA string) any {
	return gitea.UpdateFileOptions{
		FileOptions: gitea.FileOptions{
			Message:    time.Now().String(),
			BranchName: config.Data.Gitea.Branch,
			Author: gitea.Identity{
				Name: "consul-control",
			},
			Committer: gitea.Identity{
				Name: "consul-control",
			},
			Signoff: true,
		},
		SHA:      SHA,
		Content:  base64.StdEncoding.EncodeToString(content),
		FromPath: "",
	}
}

// no use ,just for interface
func (g *Opt) NewCreateCommitOptions() any {
	return nil
}

// no use ,just for interface
func (g *Opt) NewCommitActionOptions(content []byte, key string, action gitlabSDK.FileActionValue) any {
	return nil
}

func (g *Opt) PushCommit(commitOpt any) error {
	return nil
}

func (g *Opt) GetRepo() models.Repo {
	return models.GITEA
}
