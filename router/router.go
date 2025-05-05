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
			adminRouter.POST("", ctr.Admin.Login)
			adminRouter.DELETE("", ctr.Admin.Logout)
			adminRouter.GET("/", ctr.Admin.LogStatus)
			// adminRouter.Use(middleware.CheckRole(1))
			adminRouter.PUT("/", ctr.Admin.Update)
			adminRouter.GET("/stu", ctr.Admin.GetStu)
			adminRouter.PUT("/stu", ctr.Admin.UpdateStu)
			// adminRouter.Use(middleware.CheckRole(2))
			adminRouter.POST("/register", ctr.Admin.Register)
		}
		stuRouter := apiRouter.Group("/stu")
		{
			stuRouter.DELETE("", ctr.Stu.Logout)
			stuRouter.PUT("", ctr.Stu.Update)
		}
		intervRouter := apiRouter.Group("/interv")
		{
			intervRouter.POST("", ctr.Interv.New)
			intervRouter.DELETE("", ctr.Interv.Delete)
			intervRouter.PUT("", ctr.Interv.Update)
			intervRouter.GET("", ctr.Interv.Get)
		}
		queRouter := apiRouter.Group("/que")
		{
			// queRouter.Use(middleware.CheckRole(1))
			queRouter.GET("", ctr.Que.Get)
			queRouter.POST("", ctr.Que.New)
			queRouter.DELETE(":id", ctr.Que.Delete)
			queRouter.PUT("", ctr.Que.Update)
		}
	}
}
