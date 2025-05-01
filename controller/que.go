package controller

import (
	"net/http"
	"zhaoxin2025/common"
	"zhaoxin2025/model"

	"github.com/gin-gonic/gin"
)

type Que struct{}

func (*Que) Get(c *gin.Context) {
	// 传入数据结构体
	var info struct {
		NetID string `json:"netid" binding:"required,len=10,numeric"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

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

func (*Que) Delete(c *gin.Context) {
	var info struct {
		ID int `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	err := srv.Que.Delete(info.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

func (*Que) Update(c *gin.Context) {
	var info struct {
		ID   int       `json:"id" binding:"required"`
		Data model.Que `json:"data" binding:"required,dive"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	err := srv.Que.Update(info.ID, info.Data)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}
