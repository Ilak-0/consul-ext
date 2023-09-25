package cron

import (
	"consul-ext/service"
	"github.com/robfig/cron"
	"log"
)

func SyncConsulExporterTags() {
	c := cron.New()
	c.AddFunc("0 1 * * *", func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("cron sync consul exporter failed", err)
			}
		}()
		err := service.DailySyncConsulSvcs()
		if err != nil {
			log.Println("daily sync consul svcs failed:", err)
		}
	})
	c.Start()
	select {}
}
