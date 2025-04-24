package main

import (
	"fmt"
	"zhaoxin2025/config"
	"zhaoxin2025/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(config.Config.AppMode)
	srv := router.NewServer()

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("fail to init server: %s\n", err.Error())
		panic(err)
	}
}
