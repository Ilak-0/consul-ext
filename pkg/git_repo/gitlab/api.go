package gitlab

import (
	"consul-ext/pkg/config"
	"consul-ext/pkg/models"
	gitlabSDK "github.com/xanzy/go-gitlab"
	"time"
)

func (g *Opt) GetFile(key string) (*models.FileData, error) {
	g.Lock.RLock()
	defer g.Lock.RUnlock()
	opt := &gitlabSDK.GetFileOptions{
		Ref: gitlabSDK.String(config.Data.Gitlab.Branch),
	}
	file, _, err := g.Client.RepositoryFiles.GetFile(config.Data.Gitlab.Repo, key, opt)
	if err != nil {
		return nil, err
	}
	return &models.FileData{
		Content: file.Content,
		SHA:     file.LastCommitID,
	}, nil
}

func (g *Opt) UpdateFile(key string, updateFileOptions any) error {
	g.Lock.Lock()
	defer g.Lock.Unlock()
	options := updateFileOptions.(gitlabSDK.UpdateFileOptions)
	_, _, err := g.Client.RepositoryFiles.UpdateFile(config.Data.Gitlab.Repo, key, &options, nil)
	if err != nil {
		return err
	}
	return nil
}

func (g *Opt) CreateFile(key string, createFileOptions any) error {
	g.Lock.Lock()
	defer g.Lock.Unlock()
	options := createFileOptions.(gitlabSDK.CreateFileOptions)
	_, _, err := g.Client.RepositoryFiles.CreateFile(config.Data.Gitlab.Repo, key, &options, nil)
	if err != nil {
		return err
	}
	return nil
}

func (g *Opt) PushCommit(commitOpt any) error {
	commit := commitOpt.(gitlabSDK.CreateCommitOptions)
	if len(commit.Actions) == 0 {
		return nil
	}
	_, _, err := g.Client.Commits.CreateCommit(config.Data.Gitlab.Repo, &commit, nil)
	if err != nil {
		return err
	}
	return nil
}

func (g *Opt) NewEvent() any {
	return new(gitlabSDK.PushEvent)
}

func (g *Opt) NewCreateFileOptions(content []byte) any {
	return gitlabSDK.CreateFileOptions{
		Branch:        gitlabSDK.String(config.Data.Gitlab.Branch),
		Content:       gitlabSDK.String(string(content)),
		CommitMessage: gitlabSDK.String(time.Now().String()),
		AuthorName:    gitlabSDK.String("consul-control"),
	}
}

func (g *Opt) NewUpdateFileOptions(content []byte, SHA string) any {
	return gitlabSDK.UpdateFileOptions{
		Branch:        gitlabSDK.String(config.Data.Gitlab.Branch),
		AuthorName:    gitlabSDK.String("consul-control"),
		Content:       gitlabSDK.String(string(content)),
		CommitMessage: gitlabSDK.String(time.Now().String()),
		LastCommitID:  gitlabSDK.String(SHA),
	}
}

func (g *Opt) NewCreateCommitOptions() any {
	return gitlabSDK.CreateCommitOptions{
		Branch:        gitlabSDK.String(config.Data.Gitlab.Branch),
		CommitMessage: gitlabSDK.String(time.Now().String()),
		//StartSHA:      gitlabSDK.String(SHA),
		AuthorName: gitlabSDK.String("consul-control"),
	}
}

func (g *Opt) NewCommitActionOptions(content []byte, key string, action gitlabSDK.FileActionValue) any {
	return gitlabSDK.CommitActionOptions{
		Action:   &action,
		FilePath: gitlabSDK.String(key),
		Content:  gitlabSDK.String(string(content)),
	}
}

func (g *Opt) GetRepo() models.Repo {
	return models.GITLAB
}
