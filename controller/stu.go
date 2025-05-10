package controller

import (
	"net/http"
	"zhaoxin2025/common"

	"github.com/gin-gonic/gin"
)

type Stu struct{}

// 学生登录
func (*Stu) Login(c *gin.Context) {
	var info struct {
		NetID string `json:"netid" binding:"required,len=10,numeric"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	openid, err := srv.Stu.Login(info.NetID, info.Code)
	if err != nil {
		c.Error(err)
		return
	}
	SessionSet(c, "user-session", UserSession{
		ID:       NetID(info.NetID),
		Username: openid,
		Level:    1,
	})
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 学生登录状态
func (*Stu) LogStatus(c *gin.Context) {
	session := SessionGet(c, "user-session")
	if session == nil {
		c.Error(common.ErrNew(nil, common.AuthErr))
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, session))
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
	var info struct {
		NetID    string `json:"netid" binding:"required,len=10,numeric"`
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone" binding:"required,len=11,numeric"`
		Mail     string `json:"mail" binding:"required,email"`
		School   string `json:"school" binding:"required"`
		Mastered string `json:"mastered" binding:"required"`
		ToMaster string `json:"tomaster" binding:"required"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	if err := srv.Stu.Update(Struct2Map(info)); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}
