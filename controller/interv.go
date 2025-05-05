package controller

import (
	"net/http"
	"time"
	"zhaoxin2025/common"
	"zhaoxin2025/service"

	"github.com/gin-gonic/gin"
)

type Interv struct{}

// 新建面试时间
func (*Interv) New(c *gin.Context) {
	var info struct {
		Times []time.Time `json:"times" binding:"required"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	data, err := srv.Interv.New(info.Times)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, data))
}

// 更新面试记录
func (*Interv) Update(c *gin.Context) {
	var info service.IntervUpdate
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	if err := srv.Interv.Update(info); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

// 查询面试记录
func (*Interv) Get(c *gin.Context) {
	var info service.SelectStu
	if err := c.ShouldBindQuery(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	data, err := srv.Interv.Get(info)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, data)
}

// 删除面试记录
func (*Interv) Delete(c *gin.Context) {
	var info struct {
		ID []int `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	data, err := srv.Interv.Delete(info.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, data)
}
