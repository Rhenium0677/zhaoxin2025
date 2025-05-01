package controller

import (
	"net/http"
	"zhaoxin2025/common"

	"github.com/gin-gonic/gin"
)

type Stu struct{}

func (*Stu) Login(c *gin.Context) {
	var info struct {
		NetID string `json:"netid" binding:"required,len=10,numeric"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}
