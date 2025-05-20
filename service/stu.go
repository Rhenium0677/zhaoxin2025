package service

import (
	"errors"
	"time"
	"zhaoxin2025/common"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"

	"gorm.io/gorm"
)

// Stu 定义了学生相关的服务操作
type Stu struct{}

// Login 处理学生登录逻辑
// 通过微信code获取openid，查询或创建学生记录，并校验openid
func (*Stu) Login(netid string, code string) (string, error) {
	// 调用微信登录接口获取用户信息
	// _, openid, err := WxLogin(code)
	openid := "just for test remember to modify these lines"
	// if err != nil {
	// 	return "", common.ErrNew(err, common.AuthErr)
	// }
	// if openid == "" {
	// 	return "", common.ErrNew(errors.New("获取openid失败"), common.AuthErr)
	// }
	if first, err := CheckFirst(openid); err != nil {
		return "", err
	} else if first {
		// 如果是第一次登录，创建学生记录
		if err := model.DB.Model(&model.Stu{}).Create(&model.Stu{
			NetID:  netid,
			OpenID: openid,
		}).Error; err != nil {
			logger.DatabaseLogger.Errorf("创建学生记录失败: %v", err)
			return "", common.ErrNew(err, common.SysErr)
		}
		return openid, nil
	}
	var info model.Stu
	// 根据netid查询学生信息
	if err := model.DB.Model(&model.Stu{}).Where("netid = ?", netid).First(&info).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", common.ErrNew(errors.New("没找到"), common.NotFoundErr)
		}
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
func (*Stu) Update(info model.Stu) error {
	// 根据netid更新学生信息
	if err := model.DB.Model(&model.Stu{}).Where("id = ?", info.ID).Updates(&info).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新学生信息失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 更新成功
	return nil
}

func (*Stu) UpdateMessage(netid string, message int) error {
	// 更新订阅消息设置
	if err := model.DB.Model(&model.Stu{}).Where("netid = ?", netid).Update("message", message).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新订阅消息设置失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 更新成功
	return nil
}

// GetInterv 查询学生的面试记录
func (*Stu) GetInterv(netid string) (model.Interv, error) {
	// 查询学生的面试记录
	var data model.Interv
	if err := model.DB.Where("netid = ?", netid).First(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询学生面试记录失败: %v", err)
		return model.Interv{}, common.ErrNew(err, common.SysErr)
	}
	// 返回查询到的面试记录
	return data, nil
}

// AppointInterv 更新学生的面试记录
func (*Stu) AppointInterv(netid string, intervid int) error {
	// 查询面试记录
	var record model.Interv
	if err := model.DB.Model(&model.Interv{}).Where("id = ?", intervid).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("该面试不存在"), common.NotFoundErr)
		}
		logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 检查面试记录是否已经被预约
	if record.NetID != nil && *record.NetID != "" {
		return common.ErrNew(errors.New("该面试记录已经被预约"), common.AuthErr)
	}
	// 更新学生的面试记录
	if err := model.DB.Model(&model.Interv{}).Where("id = ?", intervid).Update("netid", netid).Error; err != nil {
		logger.DatabaseLogger.Errorf("预约面试失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 更新成功
	return nil
}

// CancelInterv 取消学生的面试预约
func (*Stu) CancelInterv(netid string, intervid int) error {
	// 查询面试记录
	var record model.Interv
	if err := model.DB.Model(&model.Interv{}).Where("id = ?", intervid).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("没找到"), common.NotFoundErr)
		}
		logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}

	// 检查面试记录是否属于该学生
	if record.NetID == nil || *record.NetID != netid {
		return common.ErrNew(errors.New("该面试记录不属于该学生"), common.AuthErr)
	}

	// 检查面试时间是否在半小时内或已经错过
	if record.Time.Before(time.Now().Add(30*time.Minute)) || record.Time.Before(time.Now()) {
		return common.ErrNew(errors.New("面试时间在半小时内或已经错过，无法取消预约"), common.AuthErr)
	}

	// 取消学生的面试预约
	if err := model.DB.Model(&model.Interv{}).Where("id = ?", intervid).Update("netid", nil).Error; err != nil {
		logger.DatabaseLogger.Errorf("取消预约面试失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 更新成功
	return nil
}

// CheckFirst 检查学生是否第一次登录
func CheckFirst(openid string) (bool, error) {
	// 查询学生信息
	var data model.Stu
	if err := model.DB.Model(&model.Stu{}).Where("openid = ?", openid).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil // 第一次登录
		}
		logger.DatabaseLogger.Errorf("查询学生信息失败: %v", err)
		return false, common.ErrNew(err, common.SysErr)
	}
	return false, nil // 不是第一次登录
}

// GetRes 获取学生的面试结果
func (*Stu) GetRes(netid string) (model.Interv, error) {
	var data model.Interv
	if err := model.DB.Model(&model.Interv{}).Where("netid = ?", netid).First(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询学生面试结果失败: %v", err)
		return model.Interv{}, common.ErrNew(err, common.SysErr)
	}
	return data, nil
}
