package controller

import (
	"net/http"
	"zhaoxin2025/common"
	"zhaoxin2025/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type Que struct{}

// Get 获取问题
// 根据提供的筛选条件（问题内容、部门、URL）和分页信息获取问题列表
func (*Que) Get(c *gin.Context) {
	var info struct {
		Question   string             `form:"question" binding:"omitempty"`
		Department []model.Department `form:"department" binding:"omitempty,dive,oneof=tech video art"`
		Url        string             `form:"url" binding:"omitempty"`
		common.PagerForm
	}
	// 绑定并验证查询参数
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 调用服务层获取问题数据
	total, data, err := srv.Que.Get(info.Question, info.Department, info.Url, info.PagerForm)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, struct {
		Total int64
		Data  []model.Que
	}{
		Total: total,
		Data:  data,
	}))
}

// New 新建问题
// 接收问题列表并创建新的问题记录
func (*Que) New(c *gin.Context) {
	var info struct {
		List []struct {
			Question   string           `json:"question" binding:"required"`
			Department model.Department `json:"department" binding:"required,oneof=tech video art"`
			Url        string           `json:"url" binding:"omitempty"`
		} `json:"list" binding:"required,min=1,dive"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 将问题列表转换为 []model.Que 格式
	// 以便于数据库操作
	var ques []model.Que
	for _, queInfo := range info.List {
		var que model.Que
		err := copier.Copy(&que, &queInfo)
		if err != nil {
			c.Error(common.ErrNew(err, common.SysErr))
			return
		}
		ques = append(ques, que)
	}
	err := srv.Que.New(ques)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// Delete 删除问题
// 根据提供的问题ID列表删除相应的问题记录
func (*Que) Delete(c *gin.Context) {
	var info struct {
		IDs []int `json:"ids" binding:"required,min=1,dive,gte=1"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 调用服务层删除问题
	err := srv.Que.Delete(info.IDs)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// Update 更新问题
// 根据提供的问题信息更新指定ID的问题记录
func (*Que) Update(c *gin.Context) {
	var info struct {
		ID         int              `json:"id" binding:"required"`
		Question   string           `json:"question" binding:"omitempty"`
		Department model.Department `json:"department" binding:"omitempty,oneof=tech video art"`
		Times      int              `json:"times" binding:"omitempty"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	var que model.Que
	if err := copier.Copy(&que, &info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	err := srv.Que.Update(que)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// LuckyDog 为某个幸运儿指定问题
func (*Que) LuckyDog(c *gin.Context) {
	var info struct {
		NetID string `json:"netid" binding:"required"`
		QueID int    `json:"queid" binding:"required"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	err := srv.Que.LuckyDog(info.NetID, info.QueID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}
