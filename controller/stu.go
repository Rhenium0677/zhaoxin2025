package controller

import (
	"net/http"
	"strconv"
	"zhaoxin2025/common"
	"zhaoxin2025/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type Stu struct{}

// 学生登录
// 接收学生NetID和微信code，进行登录验证并设置session
func (*Stu) Login(c *gin.Context) {
	var info struct {
		NetID string `json:"netid" binding:"required,len=10,numeric"`
		Code  string `json:"code" binding:"required"`
	}
	// 绑定并验证请求参数
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 调用服务层处理登录逻辑
	openid, err := srv.Stu.Login(info.NetID, info.Code)
	if err != nil {
		c.Error(err)
		return
	}
	// 登录成功，设置用户session
	SessionSet(c, "user-session", UserSession{
		ID:       info.NetID,
		Username: openid,
		Level:    1, // 学生level默认为1
	})
	// 返回成功响应
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 学生登录状态
// 获取当前学生的登录session信息
func (*Stu) LogStatus(c *gin.Context) {
	// 从Gin上下文中获取用户session
	session := SessionGet(c, "user-session")
	// 如果session不存在，则表示未登录
	if session == nil {
		c.JSON(http.StatusOK, ResponseNew(c, "未登录"))
		return
	}
	// 返回session信息
	c.JSON(http.StatusOK, ResponseNew(c, session))
}

// 学生登出
// 清除当前学生的登录session
func (*Stu) Logout(c *gin.Context) {
	// 检查用户是否已登录
	if SessionGet(c, "user-session") == nil {
		c.Error(common.ErrNew(nil, common.AuthErr)) // 若未登录，则返回认证错误
		return
	}
	// 删除用户session
	SessionDelete(c, "user-session")
	// 返回成功响应
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 学生修改信息
// 接收学生信息并更新数据库中对应记录
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
	// 绑定并验证请求参数
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	var stu model.Stu
	if err := copier.Copy(&stu, &info); err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	if err := srv.Stu.Update(stu); err != nil {
		c.Error(err)
		return
	}
	// 返回成功响应
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 更新学生消息订阅状态
// 根据传入的订阅选项（字符串"true"或"false"）计算一个整数掩码，并更新学生的消息设置
func (*Stu) UpdateMessage(c *gin.Context) {
	var info struct {
		Subscribe  bool `json:"subscribe" binding:"required"`
		IntervTime bool `json:"intervtime" binding:"required"`
		IntervRes  bool `json:"intervres" binding:"required"`
	}
	// 绑定并验证请求体中的JSON数据
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 从Gin上下文中获取用户session
	session := SessionGet(c, "user-session").(UserSession)
	// 将session中的用户ID转换为字符串netid
	netid := session.ID
	var message int = 0 // 初始化消息掩码
	// 根据订阅状态设置掩码位
	if info.Subscribe {
		message += 1 // 对应二进制 ...001
	}
	if info.IntervTime {
		message += 2 // 对应二进制 ...010
	}
	if info.IntervRes {
		message += 4 // 对应二进制 ...100
	}
	// 调用服务层更新学生的消息订阅状态
	if err := srv.Stu.UpdateMessage(netid, message); err != nil {
		c.Error(err)
		return
	}
	// 返回成功响应 (如果操作成功，通常会返回一个成功的JSON响应)
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 查询学生自己的面试记录
func (*Stu) GetInterv(c *gin.Context) {
	// 从Gin上下文中获取用户session
	session := SessionGet(c, "user-session").(UserSession)
	// 获取学生的面试记录
	data, err := srv.Stu.GetInterv(session.ID)
	if err != nil {
		c.Error(err)
		return
	}
	// 返回面试记录
	c.JSON(http.StatusOK, ResponseNew(c, data))
}

// 学生预约面试
func (*Stu) AppointInterv(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 从Gin上下文中获取用户session
	session := SessionGet(c, "user-session").(UserSession)
	// 调用服务层预约面试
	if err := srv.Stu.AppointInterv(session.ID, id); err != nil {
		c.Error(err)
		return
	}
	// 返回成功响应
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}
