package gitea

import (
	giteaSDK "code.gitea.io/sdk/gitea"
	"consul-ext/pkg/config"
	"fmt"
	"log"
	"sync"
)

func Init() {
	var err error
	if config.Data.Gitea.Url != "" {
		Control.Client, err = giteaSDK.NewClient(config.Data.Gitea.Url, giteaSDK.SetToken(config.Data.Gitea.Token))
	} else {
		Control.Client, err = giteaSDK.NewClient("http://"+config.Data.Gitea.Host+":"+fmt.Sprint(config.Data.Gitea.Port), giteaSDK.SetToken(config.Data.Gitea.Token))
		if err != nil {
			log.Fatalf("Failed to create gitea client: %v", err)
		}
	}
}

type Opt struct {
	Lock   sync.RWMutex
	Client *giteaSDK.Client
}

var Control = &Opt{
	Lock: sync.RWMutex{},
}
