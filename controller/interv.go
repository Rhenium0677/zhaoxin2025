package controller

import (
	"net/http"
	"time"
	"zhaoxin2025/common"
	"zhaoxin2025/model"

	"github.com/gin-gonic/gin"
)

type Interv struct{}

// 查询面试记录
func (*Interv) Get(c *gin.Context) {
	// 绑定查询参数
	var info struct {
		ID          int              `form:"id" binding:"omitempty"`
		Department  model.Department `form:"department,omitempty"`
		Interviewer string           `form:"interviewer,omitempty"`
		Pass        string           `form:"pass" binding:"omitempty,oneof=true false"`
		Date        time.Time        `form:"date" binding:"omitempty"`
		common.PagerForm
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 调用服务层获取数据
	data, err := srv.Interv.Get(Struct2Map(info))
	if err != nil {
		c.Error(err)
		return
	}
	// 返回查询结果
	c.JSON(http.StatusOK, data)
}

// 新建面试时间
func (*Interv) New(c *gin.Context) {
	// 绑定请求体
	var info struct {
		Times []time.Time `json:"times" binding:"required"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 创建面试时间
	data, err := srv.Interv.New(info.Times)
	if err != nil {
		c.Error(err)
		return
	}
	// 返回冲突的时间
	c.JSON(http.StatusOK, ResponseNew(c, data))
}

// 更新面试记录
// 更新评价、是否通过、星级等信息
func (*Interv) Update(c *gin.Context) {
	// 绑定请求体
	var info struct {
		NetID       string           `json:"netid" binding:"required"`
		Interviewer string           `json:"interviewer" binding:"omitempty"`
		Time        time.Time        `json:"time" binding:"omitempty"`
		Department  model.Department `json:"department" binding:"omitempty,oneof=tech video art"`
		Star        int              `json:"star" binding:"omitempty"`
		Evaluation  string           `json:"evaluation" binding:"omitempty"`
		Pass        string           `json:"pass" binding:"omitempty,oneof=true false"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 更新面试记录
	if err := srv.Interv.Update(Struct2Map(info)); err != nil {
		c.Error(err)
		return
	}
	// 响应
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 删除面试记录
func (*Interv) Delete(c *gin.Context) {
	// 绑定请求体
	var info struct {
		ID []int `json:"id" binding:"required"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 删除面试记录
	data, err := srv.Interv.Delete(info.ID)
	if err != nil {
		c.Error(err)
		return
	}
	// 返回删除失败的ID
	c.JSON(http.StatusOK, data)
}
