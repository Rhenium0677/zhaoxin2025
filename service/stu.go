package service

import (
	"errors"
	"zhaoxin2025/common"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"
)

// Stu 定义了学生相关的服务操作
type Stu struct{}

// Login 处理学生登录逻辑
// 通过微信code获取openid，查询或创建学生记录，并校验openid
func (*Stu) Login(netid string, code string) (string, error) {
	// 调用微信登录接口获取用户信息
	_, openid, _ := WxLogin(code) // 假设 WxLogin 是一个存在的函数
	var info model.Stu
	// 根据netid查询学生信息，如果不存在则创建该学生记录
	if err := model.DB.Where("netid = ?", netid).FirstOrCreate(&info).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询学生信息失败: %v", err)
		return "", common.ErrNew(err, common.SysErr)
	}
	// 校验数据库中存储的openid与本次登录获取的openid是否一致
	if info.OpenID != openid {
		return "", common.ErrNew(errors.New("数据库中openid与用code获取到的openid不一致"), common.AuthErr)
	}
	// 登录成功，返回openid
	return openid, nil
}

// Update 更新学生信息
// 根据传入的netid，更新学生表中的对应记录
func (*Stu) Update(info map[string]interface{}) error {
	// 从传入的map中获取netid
	netid := info["netid"].(string)
	// 根据netid更新学生信息
	if err := model.DB.Model(&model.Stu{}).Where("netid = ?", netid).Updates(info).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新学生信息失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 更新成功
	return nil
}
