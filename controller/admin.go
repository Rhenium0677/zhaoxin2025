package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"zhaoxin2025/common"
	"zhaoxin2025/model"
	"zhaoxin2025/service"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type Admin struct{}

// 无需权限

// Login 管理员的登录
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
		NetID:    data.NetID,
		Username: data.Name,
		Level:    Level(data.Level),
	})
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// Logout 管理员的注销
func (*Admin) Logout(c *gin.Context) {
	// 检查是否未登录
	session := SessionGet(c, "user-session")
	if session == nil {
		c.Error(common.ErrNew(errors.New("未登录"), common.OpErr))
		return
	}
	// 注销登录
	SessionDelete(c, "user-session")
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 权限2

// Update 管理员的更新
func (*Admin) Update(c *gin.Context) {
	var info struct {
		NetID    string `json:"netid" binding:"required,len=10,numeric"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBind(&info); err != nil {
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

// GetStu 筛选并获取学生信息
func (*Admin) GetStu(c *gin.Context) {
	var info struct {
		NetID       string           `json:"netid" binding:"omitempty,len=10,numeric"`
		Name        string           `json:"name" binding:"omitempty"`
		Phone       string           `json:"phone" binding:"omitempty,len=11,numeric"`
		School      string           `json:"school" binding:"omitempty"`
		First       model.Department `json:"first" binding:"omitempty,oneof=tech video art none"`
		Second      model.Department `json:"second" binding:"omitempty,oneof=tech video art none"`
		Pass        int              `json:"pass" binding:"omitempty,oneof=0 1"`
		Interviewer string           `json:"interviewer" binding:"omitempty"`
		Star        int              `json:"star" binding:"omitempty"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	var stu model.Stu
	if err := copier.Copy(&stu, &info); err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	var interv model.Interv
	if err := copier.Copy(&interv, &info); err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	// 获取学生信息
	data, err := srv.Admin.GetStu(stu, interv)
	if err != nil {
		c.Error(err)
		return
	}
	// 响应
	c.JSON(http.StatusOK, ResponseNew(c, data))
}

// UpdateStu 更新一个学生信息
func (*Admin) UpdateStu(c *gin.Context) {
	var info struct {
		NetID    string           `json:"netid" binding:"required,len=10,numeric"`
		Name     string           `json:"name" binding:"omitempty"`
		Phone    string           `json:"phone" binding:"omitempty"`
		School   string           `json:"school" binding:"omitempty"`
		Mastered string           `json:"mastered" binding:"omitempty"`
		ToMaster string           `json:"tomaster" binding:"omitempty"`
		First    model.Department `json:"first" binding:"omitempty,oneof=tech video art"`
		Second   model.Department `json:"second" binding:"omitempty,oneof=tech video art"`
		QueID    int              `json:"que_id" binding:"omitempty,numeric"`
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	var stu model.Stu
	if err := copier.Copy(&stu, &info); err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	stu.ID = int64(id)
	var interv model.Interv
	if err := copier.Copy(&interv, &info); err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	// 更新学生信息
	if err := srv.Admin.UpdateStu(stu); err != nil {
		c.Error(err)
		return
	}
	// 响应
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// Excelize 将学生数据输出成excel
func (*Admin) Excelize(c *gin.Context) {
	// 获取学生信息
	if err := srv.Admin.Excelize(); err != nil {
		c.Error(err)
		return
	}
	const fileName = "tenzor2025.xlsx"
	// 响应
	if _, err := http.Dir(".").Open(fileName); err != nil {
		c.Error(common.ErrNew(errors.New("指定文件不存在或无法访问: "+fileName), common.NotFoundErr))
		return
	}

	// 设置 Content-Disposition 头，确保浏览器下载文件并显示指定的文件名
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	// 设置 Content-Type 为 XLSX 文档的正确 MIME 类型
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// 发送文件
	c.File(fileName)
}

// Stat 统计学生信息并输出数据
func (*Admin) Stat(c *gin.Context) {
	// 获取统计数据
	data, err := srv.Admin.Stat()
	if err != nil {
		c.Error(err)
		return
	}
	// 响应
	c.JSON(http.StatusOK, ResponseNew(c, data))
}

// 权限3

// Register 管理员的注册
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
		NetID:    info.NetID,
		Username: info.Name,
		Level:    Level(info.Level),
	})
	// 响应
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// SetTime 设置可查询面试结果的时间
func (*Admin) SetTime(c *gin.Context) {
	var info struct {
		Time time.Time `json:"time" binding:"required"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	service.AvailableTime = info.Time
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

func (*Admin) SendResultMessage(c *gin.Context) {
	if err := srv.Admin.SendResultMessage(); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 这里是向阿里云发送请求接口?
// 原来的代码在这里收集了管理员的netid,不太理解
func (a *Admin) AliyunSendMsg(c *gin.Context) {
	data, err := srv.Admin.AliyunSendItvResMsg()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, data))
}
