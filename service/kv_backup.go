package service

import (
	"bytes"
	"consul-ext/pkg/config"
	"consul-ext/pkg/tool"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"consul-ext/pkg/consul"
	"consul-ext/pkg/git_repo"
	"consul-ext/pkg/models"
	gitlabSDK "github.com/xanzy/go-gitlab"
)

func BackupAllConsulKv() error {
	// Retrieve keys from Consul
	for _, consulAddress := range consul.AddressList {
		consulOpt, err := consul.ClientMap.GetConfig(consulAddress)
		if err != nil {
			fmt.Println("Error retrieving Consul client:", err)
			return err
		}
		keys, err := consul.GetConsulKeys(consulOpt.Client, "", "")
		if err != nil {
			fmt.Println("Error retrieving Consul keys:", err)
			return err
		}
		fileMap := make(map[string][]byte)
		// Fetch values from Consul and build the result string
		for _, key := range keys {
			value, err := consul.GetConsulValue(consulOpt.Client, key)
			if err != nil {
				fmt.Printf("Error retrieving value for key %s: %v\n", key, err)
				continue
			}
			if value == "" {
				value = "null"
			} else {
				fileMap[key] = []byte(value)
			}
		}
		switch git_repo.APIFunc().GetRepo() {
		case models.GITEA:
			err = giteaFIleToRepo(git_repo.APIFunc(), fileMap, consulAddress)
		case models.GITLAB:
			err = gitLabFileToCommit(git_repo.APIFunc(), fileMap, consulAddress)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func giteaFIleToRepo(git git_repo.GitFunc, fileMap map[string][]byte, consulAddress string) error {
	for key, value := range fileMap {
		key = consulAddress + "/" + key
		err := saveFile(git, key, value)
		if err != nil {
			fmt.Println("Error saving file:", err)
			return err
		}
	}
	return nil
}

func gitLabFileToCommit(git git_repo.GitFunc, fileMap map[string][]byte, consulAddress string) error {
	commit := git.NewCreateCommitOptions().(gitlabSDK.CreateCommitOptions)
	num := len(fileMap)
	log.Println(num)
	limitChan := make(chan bool, 200)
	errChan := make(chan error, 200)
	defer func() {
		close(limitChan)
		close(errChan)
	}()
	wg := sync.WaitGroup{}
	for key, value := range fileMap {
		wg.Add(1)
		limitChan <- true
		go func(key string, value []byte) {
			defer func() {
				wg.Done()
				<-limitChan
			}()
			log.Println(num)
			num--
			key = consulAddress + "/" + key
			data, err := git.GetFile(key)
			if err != nil && !strings.Contains(err.Error(), "404") &&
				!strings.Contains(err.Error(), "GetContentsOrList") {
				log.Println("Error retrieving file content from Gitea,key:", key, err)
				errChan <- err
			}
			if data == nil {
				createFileOptions := git_repo.APIFunc().NewCommitActionOptions(value, key, gitlabSDK.FileCreate).(gitlabSDK.CommitActionOptions)
				commit.Actions = append(commit.Actions, &createFileOptions)
			} else {
				decodedData, err := base64.StdEncoding.DecodeString(data.Content)
				if err != nil {
					log.Println("Error decoding file content：", err)
					errChan <- err
				}
				if !bytes.Equal(decodedData, value) {
					updateFileOptions := git_repo.APIFunc().NewCommitActionOptions(value, key, gitlabSDK.FileUpdate).(gitlabSDK.CommitActionOptions)
					commit.Actions = append(commit.Actions, &updateFileOptions)
				}
			}
		}(key, value)
	}
	wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
		err := git.PushCommit(commit)
		if err != nil {
			return err
		}
	}
	return nil
}

func KvWatchBackup() {
	WatchKVChangeRecord(context.Background(), config.Data.KVPefix)
}

func WatchKVChangeRecord(ctx context.Context, path string) {
	for _, consulAddress := range consul.AddressList {
		go func(consulAddress string) {
			consulIp := strings.Split(consulAddress, ":")[0]
			c, err := consul.ClientMap.GetConfig(consulAddress)
			if err != nil {
				return
			}
			watch := consul.NewWatch(c.Client, 4*time.Second, 5*time.Second)
			KVPairs, err := watch.WatchTree(ctx, path)
			if err != nil {
				log.Println(err)
				return
			}
			for {
				for _, kv := range <-KVPairs {
					if tool.PrefixStr(config.Data.ExcludeKey, kv.Key) {
						continue
					}
					err := saveFile(git_repo.APIFunc(), consulIp+"/"+kv.Key, kv.Value)
					if err != nil {
						fmt.Println("Error saving file:", err)
						return
					}
				}
			}
		}(consulAddress)
	}
}

func saveFile(git git_repo.GitFunc, key string, value []byte) error {
	data, err := git.GetFile(key)
	if err != nil && !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "GetContentsOrList") {
		log.Println("Error retrieving file content from Gitea,key:", key, err)
		return err
	}
	if data == nil {
		createFileOptions := git_repo.APIFunc().NewCreateFileOptions(value)
		err := git.CreateFile(key, createFileOptions)
		if err != nil {
			log.Println("Error creating file content in Gitea:", key, string(value), err)
			return err
		}
	} else {
		decodedData, err := base64.StdEncoding.DecodeString(data.Content)
		if err != nil {
			log.Println("Error decoding file content：", err)
			return err
		}
		if !bytes.Equal(decodedData, value) {
			updateFileOptions := git_repo.APIFunc().NewUpdateFileOptions(value, data.SHA)
			err = git.UpdateFile(key, updateFileOptions)
			if err != nil {
				fmt.Println("Error updating file content in Gitea:", key, string(value), err)
				return err
			}
		}
	}
	return nil
}
