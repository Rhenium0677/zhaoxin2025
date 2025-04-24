package router

import (
	"zhaoxin2025/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Use(middleware.Error)
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	apiRouter := r.Group("/api")
	{
		adminRouter := apiRouter.Group("/admin")
		{
			adminRouter.POST("/login", ctr.Admin.Login)
			adminRouter.DELETE("/logout", ctr.Admin.Logout)
			adminRouter.GET("/logstatus", ctr.Admin.LogStatus)
			adminRouter.Use(middleware.CheckRole(2))
			adminRouter.POST("/register", ctr.Admin.Register)
		}
	}
}
