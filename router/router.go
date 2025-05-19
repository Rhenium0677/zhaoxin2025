package router

import (
	"zhaoxin2025/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Use(middleware.Error)
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	apiRouter := r.Group("/api")
	apiRouter.GET("/session", ctr.RefreshSession, middleware.CheckRole(1))
	{
		adminRouter := apiRouter.Group("/admin")
		{
			adminRouter.POST("/", ctr.Admin.Login)
			adminRouter.DELETE("/", ctr.Admin.Logout)
			adminRouter.GET("/", ctr.Admin.LogStatus)

			// adminRouter.Use(middleware.CheckRole(2))
			adminRouter.PUT("/", ctr.Admin.Update)
			adminRouter.GET("/stu", ctr.Admin.GetStu)
			adminRouter.PUT("/stu", ctr.Admin.UpdateStu)
			adminRouter.GET("/excel", ctr.Admin.Excelize)

			// adminRouter.Use(middleware.CheckRole(3))
			adminRouter.POST("/register", ctr.Admin.Register)
		}
		stuRouter := apiRouter.Group("/stu")
		{
			stuRouter.POST("/", ctr.Stu.Login)
			stuRouter.GET("/", ctr.Stu.LogStatus)

			stuRouter.Use(middleware.CheckRole(1))
			stuRouter.DELETE("/", ctr.Stu.Logout)
			stuRouter.PUT("/", ctr.Stu.Update)
			stuRouter.PUT("/message", ctr.Stu.UpdateMessage)
			stuRouter.GET("/interv", ctr.Stu.GetInterv)
			stuRouter.PUT("/interv:id", ctr.Stu.AppointInterv)
			stuRouter.DELETE("/interv:id", ctr.Stu.CancelInterv)
		}
		intervRouter := apiRouter.Group("/interv")
		{
			intervRouter.GET("/", ctr.Interv.Get)
			intervRouter.POST("/", ctr.Interv.New)
			intervRouter.DELETE("/", ctr.Interv.Delete)
			intervRouter.PUT("/", ctr.Interv.Update)
			intervRouter.PUT("/block", ctr.Interv.BlockAndRecover)
		}
		queRouter := apiRouter.Group("/que")
		{
			queRouter.GET("/", ctr.Que.Get)

			// queRouter.Use(middleware.CheckRole(1))
			queRouter.POST("/", ctr.Que.New)
			queRouter.DELETE("/", ctr.Que.Delete)
			queRouter.PUT("/", ctr.Que.Update)
			queRouter.PUT("/lucky", ctr.Que.LuckyDog)
		}
	}
}
