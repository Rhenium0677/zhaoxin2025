package service

import (
	"errors"
	"zhaoxin2025/common"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"

	"gorm.io/gorm"
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
		return model.Admin{}, common.ErrNew(errors.New("用户名或密码错误"), common.OpErr)
	}
	return admin, nil
}

func (*Admin) Update(netid string, name string, password string) error {
	info := map[string]interface{}{
		"name":     name,
		"password": password,
	}
	// 检查管理员是否存在
	var count int64
	if err := model.DB.Model(&model.Admin{}).Where("netid = ?", netid).Count(&count).Error; err != nil {
		logger.DatabaseLogger.Errorf("检查管理员是否存在失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	if count == 0 {
		return common.ErrNew(errors.New("管理员不存在"), common.OpErr)
	}
	// 更新管理员信息
	if err := model.DB.Model(&model.Admin{}).Where("netid = ?", netid).Updates(info).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新管理员信息失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// 筛选并获取学生信息
func (*Admin) GetStu(info map[string]interface{}) ([]model.Stu, error) {
	var data []model.Stu
	db := model.DB.Model(&model.Stu{}).Where(&info)
	// 将interv相关信息提取出来
	intervInfo := make(map[string]interface{})
	for _, key := range []string{"interviewer", "star", "pass"} {
		if value, ok := info[key]; ok {
			intervInfo[key] = value
		}
	}
	if len(intervInfo) > 0 {
		// 关联查询
		db = db.Joins("JOIN intervs ON stus.netid = intervs.netid")
		for key, value := range intervInfo {
			db = db.Where("intervs."+key+" = ?", value)
		}
	}
	if err := db.Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取学生信息失败：%v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	return data, nil
}

// 更新一个学生信息
func (*Admin) UpdateStu(info map[string]interface{}) error {
	netid := info["netid"].(string)
	// 检查学生是否存在
	var count int64
	if err := model.DB.Model(&model.Stu{}).Where("netid = ?", netid).Count(&count).Error; err != nil {
		logger.DatabaseLogger.Errorf("检查学生是否存在失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	if count == 0 {
		return common.ErrNew(errors.New("学生不存在"), common.OpErr)
	}
	// 将interv相关信息提取出来以map的形式存储
	intervInfo := make(map[string]interface{})
	for _, key := range []string{"interviewer", "evaluation", "star", "pass"} {
		if value, ok := info[key]; ok {
			intervInfo[key] = value
		}
	}
	return model.DB.Transaction(func(tx *gorm.DB) error {
		// 使用事务确保操作的原子性
		// 分别更新学生信息和面试信息
		if err := tx.Model(&model.Stu{}).Where("netid = ?", netid).Updates(info).Error; err != nil {
			logger.DatabaseLogger.Errorf("更新学生信息失败：%v", err)
			return common.ErrNew(err, common.SysErr)
		}
		if err := tx.Model(&model.Interv{}).Where("netid = ?", netid).Updates(intervInfo).Error; err != nil {
			logger.DatabaseLogger.Errorf("更新学生面试信息失败：%v", err)
			return common.ErrNew(err, common.SysErr)
		}
		return nil
	})
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
