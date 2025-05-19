package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"time"
	"zhaoxin2025/common"
	"zhaoxin2025/service"
)

func CheckTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		if service.AvailableTime.After(time.Now()) {
			c.Error(common.ErrNew(errors.New("未到可查询结果的时间"), common.OpErr))
			c.Abort()
			return
		}
		c.Next()
	}
}
