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
