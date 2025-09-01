package router

import (
	"zhaoxin2025/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Use(middleware.Error)
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	r.Static("/data", "./data")
	apiRouter := r.Group("/api")
	apiRouter.GET("/", ctr.LogStatus)
	apiRouter.GET("/session", ctr.RefreshSession, middleware.CheckRole(1))
	{
		adminRouter := apiRouter.Group("/admin")
		{
			adminRouter.POST("/", ctr.Admin.Login)
			adminRouter.DELETE("/", ctr.Admin.Logout)

			adminRouter.Use(middleware.CheckRole(2))
			adminRouter.PUT("/", ctr.Admin.Update)
			adminRouter.GET("/stu", ctr.Admin.GetStu)
			adminRouter.PUT("/stu:id", ctr.Admin.UpdateStu)
			adminRouter.DELETE("/stu:id", ctr.Admin.DeleteStu)
			adminRouter.GET("/excel", ctr.Admin.Excelize)
			adminRouter.GET("/stat", ctr.Admin.Stat)
			adminRouter.GET("/aliyun", ctr.Admin.AliyunSendMsg)
			adminRouter.GET("/log", ctr.Admin.Log)
			adminRouter.GET("/getlog", ctr.Admin.DownloadLog)

			adminRouter.Use(middleware.CheckRole(3))
			adminRouter.POST("/register", ctr.Admin.Register)
			adminRouter.POST("/settime", ctr.Admin.SetTime)
			adminRouter.GET("/send", ctr.Admin.SendResultMessage)
		}
		stuRouter := apiRouter.Group("/stu")
		{
			stuRouter.POST("/", ctr.Stu.Login)
			stuRouter.DELETE("/", ctr.Stu.Logout)

			stuRouter.Use(middleware.CheckRole(1))
			stuRouter.PUT("/", ctr.Stu.Update)
			stuRouter.PUT("/message", ctr.Stu.UpdateMessage)
			stuRouter.GET("/date", ctr.Stu.GetIntervDate)
			stuRouter.GET("/interv", ctr.Stu.GetInterv)
			stuRouter.PUT("/interv:id", ctr.Stu.ReAppointInterv)
			stuRouter.POST("/interv:id", ctr.Stu.AppointInterv)
			stuRouter.DELETE("/interv:id", ctr.Stu.CancelInterv)
			stuRouter.GET("/result", ctr.Stu.GetRes, middleware.CheckTime())
		}
		intervRouter := apiRouter.Group("/interv")
		{
			intervRouter.GET("/que", ctr.Interv.GetQue)

			intervRouter.Use(middleware.CheckRole(2))
			intervRouter.POST("/", ctr.Interv.New)
			intervRouter.GET("/", ctr.Interv.Get)
			intervRouter.GET("/date", ctr.Interv.GetDate)
			intervRouter.POST("/create", ctr.Interv.Create)
			intervRouter.PUT("/cancel", ctr.Interv.Cancel)
			intervRouter.DELETE("/", ctr.Interv.Delete)
			intervRouter.PUT("/swap", ctr.Interv.Swap)
			intervRouter.PUT("/", ctr.Interv.Update)
			intervRouter.PUT("/block", ctr.Interv.BlockAndRecover)
			intervRouter.PUT("/group", ctr.Interv.QQGroup)
		}
		queRouter := apiRouter.Group("/que")
		{
			queRouter.GET("/", ctr.Que.Get)
			queRouter.GET("/:id", ctr.Que.GetOne)

			queRouter.Use(middleware.CheckRole(2))
			queRouter.POST("/", ctr.Que.New)
			queRouter.POST("/data", ctr.Que.NewData)
			queRouter.DELETE("/", ctr.Que.Delete)
			queRouter.PUT("/", ctr.Que.Update)
			queRouter.PUT("/lucky", ctr.Que.LuckyDog)
		}
	}
}
