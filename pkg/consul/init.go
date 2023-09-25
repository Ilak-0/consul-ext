package consul

import (
	"consul-ext/pkg/config"
	"github.com/hashicorp/consul/api"

	"fmt"
	"log"
	"strings"
	"sync"
)

var (
	AddressList []string
)

func Init() {
	for _, consul := range config.Data.Consul {
		//AddressList = append(AddressList, consul.Host+":"+fmt.Sprint(consul.Port))
		AddressList = append(AddressList, consul.Host)
	}

	for _, consul := range config.Data.Consul {
		defaultConfig := api.DefaultConfig()
		addr := consul.Host + ":" + fmt.Sprint(consul.Port)
		defaultConfig.Address = addr
		client, err := api.NewClient(defaultConfig)
		if err != nil {
			log.Fatal("new consul client failed:", err)
		}
		key := strings.Split(addr, ":")[0]
		ClientMap.SetConfig(key, &Opt{
			Client:  client,
			SvcName: consul.SvcName,
		})
	}
}

var ClientMap = &ControlOpts{
	Config: make(map[string]*Opt),
	Lock:   sync.RWMutex{},
}

type ControlOpts struct {
	Config map[string]*Opt
	Lock   sync.RWMutex
}

type Opt struct {
	Client  *api.Client
	SvcName []string
}

func (c *ControlOpts) GetConfig(addr string) (*Opt, error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	if client, ok := c.Config[addr]; ok {
		return client, nil
	}
	return nil, fmt.Errorf("client not found")
}

func (c *ControlOpts) SetConfig(addr string, config *Opt) {
	c.Lock.Lock()
	defer func() {
		c.Lock.Unlock()
		if err := recover(); err != nil {
			log.Println("set consul config failed:", err)
		}
	}()
	c.Config[addr] = config
}
