package controller

import (
	"strconv"
	"zhaoxin2025/model"
)

// 将netid转换为session需要的int类型
func NetID(netid string) int {
	var convErr error
	var netidInt int
	// 传入参数时做过校验，故不应出现错误
	netidInt, convErr = strconv.Atoi(netid)
	if convErr != nil {
		return 0
	}
	return netidInt
}

// 将管理员级别super normal转换成session的role数字
func Level(level model.AdminLevel) int {
	switch level {
	case model.Super:
		return 2
	case model.Normal:
		return 1
	default:
		return 0 // 默认值，处理未知的管理员级别
	}
}
