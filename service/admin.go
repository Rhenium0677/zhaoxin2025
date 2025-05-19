package service

import (
	"errors"
	"time"
	"zhaoxin2025/common"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"

	"gorm.io/gorm"
)

var AvailableTime = time.Date(2025, time.October, 1, 0, 0, 0, 0, time.Local)

type Admin struct{}

func (*Admin) Login(netid string, password string) (model.Admin, error) {
	// 检查管理员是否存在
	var admin model.Admin
	if err := model.DB.Where("netid = ?", netid).First(&admin).Error; err != nil {
		return model.Admin{}, common.ErrNew(err, common.OpErr)
	}
	// 加密传入的密码
	encrpted, err := Encrypt(password)
	if err != nil {
		return model.Admin{}, common.ErrNew(err, common.SysErr)
	}
	// 检查密码是否正确
	if encrpted != admin.Password {
		return model.Admin{}, common.ErrNew(errors.New("用户名或密码错误"), common.OpErr)
	}
	return admin, nil
}

func (*Admin) Update(netid string, name string, password string) error {
	info := map[string]any{
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
	if err := model.DB.Model(&model.Admin{}).Where("netid = ?", netid).Updates(&info).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新管理员信息失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// 筛选并获取学生信息
func (*Admin) GetStu(stuInfo model.Stu, intervInfo model.Interv) ([]model.Stu, error) {
	var data []model.Stu
	db := model.DB.Model(&model.Stu{}).Where(&stuInfo)
	if intervInfo.Evaluation != "" {
		// 关联查询
		db = db.Joins("JOIN intervs ON stus.netid = intervs.netid").Where(&intervInfo)
	}
	if err := db.Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取学生信息失败：%v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	return data, nil
}

// 更新一个学生信息
func (*Admin) UpdateStu(stuInfo model.Stu, intervInfo model.Interv) error {
	// 检查学生是否存在
	var existed model.Stu
	if err := model.DB.Model(&model.Stu{}).Where("id = ?", stuInfo.ID).First(&existed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("学生不存在"), common.NotFoundErr)
		}
		logger.DatabaseLogger.Errorf("检查学生是否存在失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return model.DB.Transaction(func(tx *gorm.DB) error {
		// 使用事务确保操作的原子性
		// 分别更新学生信息和面试信息
		if err := tx.Model(&model.Stu{}).Where("id = ?", stuInfo.ID).Updates(stuInfo).Error; err != nil {
			logger.DatabaseLogger.Errorf("更新学生信息失败：%v", err)
			return common.ErrNew(err, common.SysErr)
		}
		if err := tx.Model(&model.Interv{}).Where("id = ?", stuInfo.ID).Updates(intervInfo).Error; err != nil {
			logger.DatabaseLogger.Errorf("更新学生面试信息失败：%v", err)
			return common.ErrNew(err, common.SysErr)
		}
		return nil
	})
}

// Excelize 获取学生信息并导出为excel
func (*Admin) Excelize() error {
	// 获取所有学生信息
	var data []model.Stu
	if err := model.DB.Model(&model.Stu{}).Preload("Interv.Que").Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取学生信息失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 创建excel文件
	if err := Excelize(data, "tenzor2025.xlsx"); err != nil {
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// Register 管理员注册
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

// SetTime 设置可查询面试结果的时间
func (*Admin) SetTime(time time.Time) {
	AvailableTime = time
}
