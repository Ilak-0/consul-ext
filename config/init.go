package config

import (
	"consul-ext/db"
	"consul-ext/pkg/config"
	"consul-ext/pkg/consul"
	"consul-ext/pkg/git_repo/gitea"
	"consul-ext/pkg/git_repo/gitlab"
	"log"
)

func Init() {
	// consul
	consul.Init()
	// gitea
	if config.Data.Gitea != nil {
		gitea.Init()
		log.Println("gitea init success:", config.Data.Gitea.Owner, config.Data.Gitea.Repo, config.Data.Gitea.Branch)
	}
	// gitlab
	if config.Data.Gitlab != nil {
		gitlab.Init()
		log.Println("gitlab init success:", config.Data.Gitlab.Host, config.Data.Gitlab.Port, config.Data.Gitlab.Repo, config.Data.Gitlab.Branch)
	}
	// db
	if config.Data.Database != nil {
		db.Init()
		log.Println("db init success")
	}
}
