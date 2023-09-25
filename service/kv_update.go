package service

import (
	"consul-ext/pkg/consul"
	"consul-ext/pkg/git_repo"
	"consul-ext/pkg/models"
	gitlabSDK "github.com/xanzy/go-gitlab"

	"strings"
)

func UpdateFromGiteaRepo(event *models.GiteaPayloadCommit) error {
	for _, commit := range event.Commits {
		for _, modified := range commit.Modified {
			err := UpdateKV(modified)
			if err != nil {
				return err
			}
		}
		for _, added := range commit.Added {
			err := createKV(added)
			if err != nil {
				return err
			}
		}
		for _, deleted := range commit.Removed {
			err := deleteKV(deleted)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func UpdateFromGitLabRepo(event *gitlabSDK.PushEvent) error {
	for _, commit := range event.Commits {
		for _, modified := range commit.Modified {
			err := UpdateKV(modified)
			if err != nil {
				return err
			}
		}
		for _, added := range commit.Added {
			err := createKV(added)
			if err != nil {
				return err
			}
		}
		for _, deleted := range commit.Removed {
			err := deleteKV(deleted)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func UpdateKV(path string) error {
	data, err := git_repo.APIFunc().GetFile(path)
	if err != nil && !strings.Contains(err.Error(), "404") {
		return err
	}
	strSlices := strings.Split(path, "/")
	consulAddr := strSlices[0]
	consulConfig, err := consul.ClientMap.GetConfig(consulAddr)
	if err != nil {
		return err
	}
	consulPath := strings.Join(strSlices[1:], "/")
	err = consul.UpdateKv(consulConfig.Client, consulPath, data.Content)
	if err != nil {
		return err
	}
	return nil
}

func createKV(path string) error {
	data, err := git_repo.APIFunc().GetFile(path)
	if err != nil && !strings.Contains(err.Error(), "404") {
		return err
	}
	strSlices := strings.Split(path, "/")
	consulAddr := strSlices[0]
	consulConfig, err := consul.ClientMap.GetConfig(consulAddr)
	if err != nil {
		return err
	}
	consulPath := strings.Join(strSlices[1:], "/")
	err = consul.UpdateKv(consulConfig.Client, consulPath, data.Content)
	if err != nil {
		return err
	}
	return nil
}

func deleteKV(path string) error {
	strSlices := strings.Split(path, "/")
	consulAddr := strSlices[0]
	consulConfig, err := consul.ClientMap.GetConfig(consulAddr)
	if err != nil {
		return err
	}
	consulPath := strings.Join(strSlices[1:], "/")
	err = consul.DeleteKv(consulConfig.Client, consulPath)
	if err != nil {
		return err
	}
	return nil
}
