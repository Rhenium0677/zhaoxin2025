package service

import (
	"errors"
	"zhaoxin2025/common"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"
)

type Admin struct{}

func (*Admin) Login(netid string, password string) (model.Admin, error) {
	// 检查管理员是否存在
	var admin model.Admin
	if err := model.DB.Where("netid = ?", netid).First(&admin).Error; err != nil {
		return model.Admin{}, common.ErrNew(err, common.OpErr)
	}
	// 解密密码
	nosecret, err := Decrypt(admin.Password)
	if err != nil {
		return model.Admin{}, common.ErrNew(err, common.SysErr)
	}
	// 检查密码是否正确
	if nosecret != password {
		return model.Admin{}, common.ErrNew(errors.New("密码错误"), common.OpErr)
	}
	return admin, nil
}

func (*Admin) Update(netid string, name string, password string) error {
	// 检查管理员是否存在
	var count int64
	if err := model.DB.Model(&model.Admin{}).Where("netid = ?", netid).Count(&count).Error; err != nil {
		logger.DatabaseLogger.Errorf("检查管理员是否存在失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	if count == 0 {
		return common.ErrNew(errors.New("管理员不存在"), common.OpErr)
	}
	data := model.Admin{
		Name:     name,
		Password: password,
	}
	// 更新管理员信息
	if err := model.DB.Model(&model.Admin{}).Where("netid = ?", netid).Updates(data).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新管理员信息失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// 获取一个学生信息
func (*Admin) GetStu(info SelectStu) ([]model.Stu, error) {
	var data []model.Stu
	if err := model.DB.Model(&model.Stu{}).Where(&info).Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取学生信息失败：%v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	return data, nil
}

// 管理员注册
func (*Admin) Register(netid string, name string, password string, level model.AdminLevel) error {
	// 检查管理员是否存在
	var count int64
	if err := model.DB.Model(&model.Admin{}).Where("netid = ?", netid).Count(&count).Error; err != nil {
		logger.DatabaseLogger.Errorf("检查管理员是否存在失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	if count > 0 {
		return common.ErrNew(errors.New("管理员已存在"), common.OpErr)
	}
	// 用data结构体存储管理员信息
	// 加密密码
	secret, err := Encrypt(password)
	if err != nil {
		return common.ErrNew(err, common.SysErr)
	}
	data := model.Admin{
		NetID:    netid,
		Name:     name,
		Password: secret,
		Level:    level,
	}
	// 在数据库中插入数据
	if err := model.DB.Model(&model.Admin{}).Create(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("插入管理员信息失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}
