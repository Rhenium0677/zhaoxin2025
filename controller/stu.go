package controller

import (
	"net/http"
	"zhaoxin2025/common"
	"zhaoxin2025/service"

	"github.com/gin-gonic/gin"
)

type Stu struct{}

// 学生登录
func (*Stu) Login(c *gin.Context) {
	var info struct {
		NetID string `json:"netid" binding:"required,len=10,numeric"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// token, openid, err := srv.Stu.Login(info.NetID, info.Code)
	// if err != nil {
	// 	c.Error(err)
	// 	return
	// }
	// SessionSet(c, "user-session", UserSession{
	// 	ID:    NetID(info.NetID),
	// 	Level: 0,
	// })
	// c.JSON(http.StatusOK, ResponseNew(c, gin.H{
	// 	"token": token,
	// }))
}

// 学生登出
func (*Stu) Logout(c *gin.Context) {
	if SessionGet(c, "user-session") == nil {
		c.Error(common.ErrNew(nil, common.AuthErr))
		return
	}
	SessionDelete(c, "user-session")
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 学生修改信息
func (*Stu) Update(c *gin.Context) {
	//
	//
	var info service.StuUpdate
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	if err := srv.Stu.Update(info); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}
