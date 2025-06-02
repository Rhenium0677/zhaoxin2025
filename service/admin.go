package service

import (
	"errors"
	"time"
	"zhaoxin2025/common"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"

	"gorm.io/gorm"
)

var AvailableTime = time.Date(2026, time.October, 1, 0, 0, 0, 0, time.Local)
var BlockTable = make(map[TimeRange]struct{})

type Admin struct{}

func (*Admin) Login(netid string, password string) (model.Admin, error) {
	// 检查管理员是否存在
	var record model.Admin
	if err := model.DB.Model(&model.Admin{}).Where("netid = ?", netid).First(&record).Error; err != nil {
		return model.Admin{}, common.ErrNew(err, common.OpErr)
	}
	// 加密传入的密码
	decrypted, err := Decrypt(record.Password)
	if err != nil {
		return model.Admin{}, common.ErrNew(err, common.SysErr)
	}
	// 检查密码是否正确
	if decrypted != password {
		return model.Admin{}, common.ErrNew(errors.New("用户名或密码错误"), common.OpErr)
	}
	return record, nil
}

func (*Admin) Update(netid string, name string, password string) error {
	info := model.Admin{
		Name:     name,
		Password: password,
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
	db := model.DB.Model(&model.Stu{}).Preload("Interv").Where(&stuInfo)
	if intervInfo.Evaluation != "" {
		// 关联查询
		db = db.Where("interv.evaluation LIKE ?", intervInfo.Evaluation)
	}
	if intervInfo.Interviewer != "" {
		// 关联查询
		db = db.Where("interv.interviewer LIKE ?", intervInfo.Interviewer)
	}
	if err := db.Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取学生信息失败：%v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	return data, nil
}

// 更新一个学生信息
func (*Admin) UpdateStu(stuInfo model.Stu) error {
	// 检查学生是否存在
	var existed model.Stu
	if err := model.DB.Model(&model.Stu{}).Where("id = ?", stuInfo.ID).First(&existed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("学生不存在"), common.NotFoundErr)
		}
		logger.DatabaseLogger.Errorf("检查学生是否存在失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	if err := model.DB.Model(&model.Stu{}).Where("id = ?", stuInfo.ID).Updates(stuInfo).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新学生信息失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// Excelize 获取学生信息并导出为excel
func (*Admin) Excelize() error {
	// 获取所有学生信息
	var data []model.Stu
	if err := model.DB.Model(&model.Stu{}).Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取学生信息失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 创建excel文件
	if err := Excelize(data, "tenzor2025.xlsx"); err != nil {
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// Stat 获取学生信息并统计
func (*Admin) Stat() (Stat, error) {
	var data []model.Stu
	if err := model.DB.Model(&model.Stu{}).Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取学生信息失败：%v", err)
		return Stat{}, common.ErrNew(err, common.SysErr)
	}
	return GetStat(data), nil
}

// Register 管理员注册
func (*Admin) Register(netid string, name string, password string, level model.AdminLevel) error {
	// 检查管理员是否存在
	var record model.Admin
	if err := model.DB.Model(&model.Admin{}).Where("netid = ?", netid).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用data结构体存储管理员信息
			// 加密密码
			encrypted, err := Encrypt(password)
			if err != nil {
				return common.ErrNew(err, common.SysErr)
			}
			data := model.Admin{
				NetID:    netid,
				Name:     name,
				Password: encrypted,
				Level:    level,
			}
			// 在数据库中插入数据
			if err := model.DB.Model(&model.Admin{}).Create(&data).Error; err != nil {
				logger.DatabaseLogger.Errorf("插入管理员信息失败：%v", err)
				return common.ErrNew(err, common.SysErr)
			}
			return nil

		}
		logger.DatabaseLogger.Errorf("检查管理员是否存在失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return common.ErrNew(errors.New("管理员已存在"), common.OpErr)
}

func (*Admin) SendResultMessage() error {
	// 获取所有学生信息
	var data []model.Stu
	if err := model.DB.Model(&model.Stu{}).Where("message > 3").Preload("Interv").Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取学生信息失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	for _, stu := range data {
		err := ResultQueue.AddMessage(stu)
		if err != nil {
			println("添加面试结果订阅消息失败: %v，学生netid: %s", err, stu.NetID)
		}
	}
	return nil
}
