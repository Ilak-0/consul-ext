package main

import (
	"consul-ext/config"
	"consul-ext/cron"
	cfg "consul-ext/pkg/config"
	"consul-ext/router"
	"consul-ext/service"
	"log"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

// @title Consul Control API
// @version 1.0
// @description A manager for hosts
// @termsOfService http://127.0.0.1:8080
// @license.name MIT
func main() {
	config.Init() // config init
	go cron.SyncConsulExporterTags()
	if cfg.Data.ConsulKVWatch {
		go service.KvWatchBackup()
	}
	router.Run()
}
