package controller

import (
	"errors"
	"net/http"
	"zhaoxin2025/common"
	"zhaoxin2025/model"
	"zhaoxin2025/service"

	"github.com/gin-gonic/gin"
)

type Admin struct{}

// 无需权限
// 管理员的登录
func (*Admin) Login(c *gin.Context) {
	// 传入数据结构体
	var info struct {
		NetID    string `json:"netid" binding:"required,len=10,numeric"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 登录
	data, err := srv.Admin.Login(info.NetID, info.Password)
	if err != nil {
		c.Error(err)
		return
	}
	// 设置session
	SessionSet(c, "user-session", UserSession{
		ID:       NetID(data.NetID),
		Username: data.Name,
		Level:    Level(data.Level),
	})
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 管理员的注销
func (*Admin) Logout(c *gin.Context) {
	// 检查是否未登录
	session := SessionGet(c, "user-session")
	if session == nil {
		c.Error(common.ErrNew(errors.New("未登录"), common.OpErr))
		return
	}
	// 注销登录
	SessionClear(c)
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 管理员的登录状态
func (*Admin) LogStatus(c *gin.Context) {
	// 获取session
	session := SessionGet(c, "user-session")
	if session == nil {
		c.JSON(http.StatusOK, ResponseNew(c, nil))
		return
	}
	// 响应
	c.JSON(http.StatusOK, ResponseNew(c, session))
}

// 权限1
// 管理员的更新
func (*Admin) Update(c *gin.Context) {
	var info struct {
		NetID    string `json:"netid" binding:"required,len=10,numeric"`
		Name     string `json:"name" binding:"omitempty"`
		Password string `json:"password" binding:"omitempty"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 更新管理员信息
	err := srv.Admin.Update(info.NetID, info.Name, info.Password)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 筛选并获取学生信息
func (*Admin) GetStu(c *gin.Context) {
	var info service.SelectStu
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 获取学生信息
	data, err := srv.Admin.GetStu(info)
	if err != nil {
		c.Error(err)
		return
	}
	// 响应
	c.JSON(http.StatusOK, ResponseNew(c, data))
}

// 更新一个学生信息
func (*Admin) UpdateStu(c *gin.Context) {
	var info service.UpdateStu
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 更新学生信息
	if err := srv.Admin.UpdateStu(info); err != nil {
		c.Error(err)
		return
	}
	// 响应
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 权限2
// 管理员的注册
func (*Admin) Register(c *gin.Context) {
	// 传入数据结构体
	var info struct {
		NetID    string           `json:"netid" binding:"required,len=10,numeric"`
		Name     string           `json:"name" binding:"required"`
		Password string           `json:"password" binding:"required"`
		Level    model.AdminLevel `json:"level" binding:"required,oneof=normal super"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 获取管理员级别
	err := srv.Admin.Register(info.NetID, info.Name, info.Password, info.Level)
	if err != nil {
		c.Error(err)
		return
	}
	// 设置session
	SessionSet(c, "user-session", UserSession{
		ID:       NetID(info.NetID),
		Username: info.Name,
		Level:    Level(info.Level),
	})
	// 响应
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}
