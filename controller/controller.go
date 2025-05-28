package controller

import (
	"errors"
	"net/http"
	"zhaoxin2025/common"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Admin
	Stu
	Interv
	Que
}

func New() *Controller {
	Controller := &Controller{}
	return Controller
}

func (*Controller) RefreshSession(c *gin.Context) {
	user, ok := SessionGet(c, "user-session").(UserSession)
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}
	SessionDelete(c, "user-session")
	SessionSet(c, "user-session", user)
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

func (*Controller) LogStatus(c *gin.Context) {
	// 从Gin上下文中获取用户session
	session := SessionGet(c, "user-session")
	// 如果session不存在，则表示未登录
	if session == nil {
		c.JSON(http.StatusOK, ResponseNew(c, "未登录"))
		return
	}
	// 返回session信息
	c.JSON(http.StatusOK, ResponseNew(c, struct {
		NetID string `json:"netid"`
	}{
		NetID: session.(UserSession).NetID,
	}))
}
