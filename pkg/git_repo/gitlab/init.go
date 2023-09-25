package gitlab

import (
	"consul-ext/pkg/config"
	"fmt"
	gitlabSDK "github.com/xanzy/go-gitlab"
	"log"
	"sync"
)

func Init() {
	var err error
	if config.Data.Gitlab.Url != "" {
		Control.Client, err = gitlabSDK.NewClient(config.Data.Gitlab.Token, gitlabSDK.WithBaseURL(config.Data.Gitlab.Url))
		if err != nil {
			log.Fatalf("Failed to create gitlab client: %v", err)
		}
	} else {
		Control.Client, err = gitlabSDK.NewClient(config.Data.Gitlab.Token, gitlabSDK.WithBaseURL("http://"+config.Data.Gitlab.Host+":"+fmt.Sprint(config.Data.Gitlab.Port)))
		if err != nil {
			log.Fatalf("Failed to create gitlab client: %v", err)
		}
	}
}

type Opt struct {
	Lock   sync.RWMutex
	Client *gitlabSDK.Client
}

var Control = &Opt{
	Lock: sync.RWMutex{},
}
