package controller

import (
	"net/http"
	"zhaoxin2025/common"
	"zhaoxin2025/model"
	"zhaoxin2025/service"

	"github.com/gin-gonic/gin"
)

type Que struct{}

// 获取问题
func (*Que) Get(c *gin.Context) {
	var info struct {
		Question   string           `form:"question" binding:"omitempty"`
		Department model.Department `form:"department" binding:"omitempty,oneof=tech video art"`
		Url        string           `form:"url" binding:"omitempty"`
		common.PagerForm
	}
	if err := c.ShouldBindQuery(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	err, data := srv.Que.Get(info.Question, info.Department, info.Url, info.PagerForm)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, data))
}

// 新建问题
func (*Que) New(c *gin.Context) {
	var info struct {
		List []model.Que `json:"list" binding:"required,dive"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	err := srv.Que.New(info.List)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 删除问题
func (*Que) Delete(c *gin.Context) {
	var info struct {
		IDs []int `json:"ids" binding:"required,dive,gte=1"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	err := srv.Que.Delete(info.IDs)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 更新问题
func (*Que) Update(c *gin.Context) {
	var info service.UpdateQue
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	err := srv.Que.Update(info)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}
