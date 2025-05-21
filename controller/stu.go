package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"zhaoxin2025/common"
	"zhaoxin2025/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type Stu struct{}

// Login 学生登录
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

// Update 学生修改信息
// 接收学生信息并更新数据库中对应记录
func (*Stu) Update(c *gin.Context) {
	var info struct {
		ID       int              `json:"id" binding:"required"`
		NetID    string           `json:"netid" binding:"required,len=10,numeric"`
		Name     string           `json:"name" binding:"required"`
		Phone    string           `json:"phone" binding:"required,len=11,numeric"`
		Mail     string           `json:"mail" binding:"required,email"`
		School   string           `json:"school" binding:"required"`
		Mastered string           `json:"mastered" binding:"required"`
		ToMaster string           `json:"tomaster" binding:"required"`
		First    model.Department `json:"first" binding:"required,oneof=tech video art none"`
		Second   model.Department `json:"second" binding:"required,oneof=tech video art none"`
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

// UpdateMessage 更新学生消息订阅状态
// 根据传入的订阅选项（字符串"true"或"false"）计算一个整数掩码，并更新学生的消息设置
func (*Stu) UpdateMessage(c *gin.Context) {
	var info struct {
		Subscribe  int `json:"subscribe" binding:"oneof=0 1"`
		IntervTime int `json:"intervtime" binding:"oneof=0 1"`
		IntervRes  int `json:"intervres" binding:"oneof=0 1"`
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
	message := 1*info.Subscribe + 2*info.IntervTime + 4*info.IntervRes
	// 调用服务层更新学生的消息订阅状态
	if err := srv.Stu.UpdateMessage(netid, message); err != nil {
		c.Error(err)
		return
	}
	// 返回成功响应 (如果操作成功，通常会返回一个成功的JSON响应)
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// GetIntervDate 查询可用的面试日期
func (*Stu) GetIntervDate(c *gin.Context) {
	data, err := srv.Stu.GetIntervDate()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, data))
}

// GetInterv 查询某天可用与不可用的面试数量
func (*Stu) GetInterv(c *gin.Context) {
	session := SessionGet(c, "user-session").(UserSession)
	var info struct {
		Date time.Time `form:"date" binding:"required"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 获取面试记录
	data, err := srv.Stu.GetInterv(info.Date)
	if err != nil {
		c.Error(err)
		return
	}
	var available, unavailable []model.Interv
	for _, value := range data {
		if value.NetID != nil && *(value.NetID) != session.ID {
			unavailable = append(unavailable, value)
		} else {
			available = append(available, value)
		}

	}
	// 返回面试记录
	c.JSON(http.StatusOK, ResponseNew(c, struct {
		Available   []model.Interv `json:"available"`
		Unavailable []model.Interv `json:"unavailable"`
	}{
		Available:   available,
		Unavailable: unavailable,
	}))
}

// AppointInterv 学生预约面试
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

// CancelInterv 学生取消预约面试
func (*Stu) CancelInterv(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 从Gin上下文中获取用户session
	session := SessionGet(c, "user-session").(UserSession)
	// 调用服务层取消预约面试
	if err := srv.Stu.CancelInterv(session.ID, id); err != nil {
		c.Error(err)
		return
	}
	// 返回成功响应
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// GetResult 查询学生面试结果
func (*Stu) GetRes(c *gin.Context) {
	session := SessionGet(c, "user-session").(UserSession)
	data, err := srv.Stu.GetRes(session.ID)
	if err != nil {
		c.Error(err)
		return
	}

	fileName := fmt.Sprintf("%s.json", data.Department)
	targetDir := "QQGroup"
	path := filepath.Join(targetDir, fileName)
	fmt.Printf("尝试在路径 %s 读取文件\n", path)

	file, err := os.Open(path)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	defer file.Close()

	var qqGroup struct {
		URL        string           `json:"url"`
		QQGroup    string           `json:"qqgroup"`
		Department model.Department `json:"department"`
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&qqGroup)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, struct {
		Data    model.Interv `json:"data"`
		URL     string       `json:"url"`
		QQGroup string       `json:"qqgroup"`
	}{
		Data:    data,
		URL:     qqGroup.URL,
		QQGroup: qqGroup.QQGroup,
	}))
}
