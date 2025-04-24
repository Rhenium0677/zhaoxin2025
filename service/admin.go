package service

import (
	"errors"
	"zhaoxin2025/common"
	"zhaoxin2025/model"
)

type Admin struct{}

func (*Admin) Login(netid string, password string) (model.Admin, error) {
	// 检查管理员是否存在
	var admin model.Admin
	if err := model.DB.Where("netid = ?", netid).First(&admin).Error; err != nil {
		return model.Admin{}, common.ErrNew(err, common.OpErr)
	}
	// 检查密码是否正确
	if admin.Password != password {
		return model.Admin{}, common.ErrNew(errors.New("密码错误"), common.OpErr)
	}
	return admin, nil
}

// 管理员注册
func (*Admin) Register(netid string, name string, password string, level model.AdminLevel) error {
	// 检查管理员是否存在
	var count int64
	if err := model.DB.Model(&model.Admin{}).Where("netid = ?", netid).Count(&count).Error; err != nil {
		return common.ErrNew(err, common.SysErr)
	}
	if count > 0 {
		return common.ErrNew(errors.New("管理员已存在"), common.OpErr)
	}
	// 用data结构体存储管理员信息
	data := model.Admin{
		NetID:    netid,
		Name:     name,
		Password: password,
		Level:    level,
	}
	// 在数据库中插入数据
	if err := model.DB.Model(&model.Admin{}).Create(&data).Error; err != nil {
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}
