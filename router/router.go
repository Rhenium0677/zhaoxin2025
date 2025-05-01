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
			// adminRouter.Use(middleware.CheckRole(1))
			adminRouter.PUT("/:id", ctr.Admin.Update)
			adminRouter.GET("/stu", ctr.Admin.GetStu)
			// adminRouter.Use(middleware.CheckRole(2))
			adminRouter.POST("/register", ctr.Admin.Register)
		}
		queRouter := apiRouter.Group("/que")
		{
			// queRouter.Use(middleware.CheckRole(1))
			queRouter.GET("", ctr.Que.Get)
			queRouter.POST("", ctr.Que.New)
			queRouter.DELETE("/delete", ctr.Que.Delete)
			queRouter.PUT("/:id", ctr.Que.Update)
		}
	}
}
