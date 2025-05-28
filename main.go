package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"zhaoxin2025/config"
	"zhaoxin2025/router"
	"zhaoxin2025/service"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(config.Config.AppMode)
	srv := router.NewServer()

	go func() {
		c := cron.New(cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger)))
		_, err := c.AddFunc("@every 1h50m", func() {
			service.RefreshAccessToken()
		})
		if err != nil {
			fmt.Printf("fail to init cron: %s\n", err.Error())
			panic(err)
		}
		c.Start()
	}()
	service.Send()

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("fail to init server: %s\n", err.Error())
		panic(err)
	}
}
