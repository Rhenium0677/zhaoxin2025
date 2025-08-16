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
func (*Stu) Login(openid string) (bool, model.Stu, error) {
	var record model.Stu
	// 根据openid查询学生信息
	if err := model.DB.Model(&model.Stu{}).Where("openid = ?", openid).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果是第一次登录，创建学生记录
			record.OpenID = openid
			record.NetID = openid
			if err := model.DB.Model(&model.Stu{}).Create(&record).Error; err != nil {
				logger.GinLogger.Errorf("创建学生记录失败: %v", err)
				return true, model.Stu{}, common.ErrNew(err, common.SysErr)
			}
			return true, record, nil
		}
		logger.GinLogger.Errorf("查询学生信息失败: %v", err)
		return false, model.Stu{}, common.ErrNew(err, common.SysErr)
	}
	// 登录成功，返回openid
	return false, record, nil
}

// Update 更新学生信息
// 根据传入的id，更新学生表中的对应记录
func (*Stu) Update(info model.Stu) error {
	var record model.Stu
	if err := model.DB.Model(&model.Stu{}).Where("id = ?", info.ID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("没找到"), common.NotFoundErr)
		}
		logger.GinLogger.Errorf("查询学生信息失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 根据id更新学生信息
	if err := model.DB.Model(&model.Stu{}).Where("id = ?", info.ID).Updates(&info).Error; err != nil {
		logger.GinLogger.Errorf("更新学生信息失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 更新成功
	return nil
}

func (*Stu) UpdateMessage(netid string, message int) error {
	var record model.Stu
	// 查询学生信息
	if err := model.DB.Model(&model.Stu{}).Where("netid = ?", netid).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("没找到"), common.NotFoundErr)
		}
		logger.GinLogger.Errorf("查询学生信息失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 更新订阅消息设置
	if err := model.DB.Model(&model.Stu{}).Where("netid = ?", netid).Update("message", message).Error; err != nil {
		logger.GinLogger.Errorf("更新订阅消息设置失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 更新成功
	return nil
}

// GetIntervDate 获取可用的面试日期
func (*Stu) GetIntervDate() (map[string]int, error) {
	// 查询可用的面试日期
	var record []model.Interv
	if err := model.DB.Model(&model.Interv{}).Where("netid IS NULL").Where("time > ?", time.Now()).Find(&record).Error; err != nil {
		logger.GinLogger.Errorf("查询可用面试日期失败: %v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	IntervDate := make(map[string]int)
	for _, interv := range record {
		IntervDate[Date(interv.Time)]++
	}
	// 返回可用的面试日期
	return IntervDate, nil
}

// GetInterv 查询某日可用面试记录
func (*Stu) GetInterv(date time.Time) ([]model.Interv, error) {
	// 查询学生的面试记录
	var data []model.Interv
	timeRange := DayRange(date)
	if err := model.DB.Where("time BETWEEN ? AND ?", timeRange.Start, timeRange.End).Where("time > ?", time.Now()).Find(&data).Error; err != nil {
		logger.GinLogger.Errorf("查询学生面试记录失败: %v", err)
		return nil, common.ErrNew(err, common.SysErr)
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
		logger.GinLogger.Errorf("查询面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 检查是否在不可预约时间段内
	for timeRange, _ := range BlockTable {
		if record.Time.After(timeRange.Start) && record.Time.Before(timeRange.End) {
			return common.ErrNew(errors.New("当前面试信息无法修改"), common.AuthErr)
		}
	}
	if time.Now().After(record.Time) {
		return common.ErrNew(errors.New("面试时间已过"), common.OpErr)
	}
	// 检查面试记录是否已经被预约
	if record.NetID != nil && *record.NetID != netid {
		return common.ErrNew(errors.New("该面试已经被预约"), common.AuthErr)
	}
	if record.NetID != nil && *record.NetID == netid {
		return nil
	}
	// 更新学生的面试记录
	var stu model.Stu
	return model.DB.Model(&model.Interv{}).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ?", intervid).Where("netid is NULL").Update("netid", netid)
		if result.Error != nil {
			logger.GinLogger.Errorf("预约面试失败: %v", result.Error)
			return common.ErrNew(result.Error, common.SysErr)
		}
		if result.RowsAffected == 0 {
			return common.ErrNew(errors.New("预约面试失败，请稍后再试"), common.SysErr)
		}
		if stu.Message%2 == 1 {
			err := SendRegister(stu)
			logger.GinLogger.Infof("尝试添加订阅消息, openid: %s", stu.OpenID)
			if err != nil {
				logger.GinLogger.Errorf("添加订阅消息失败: %v", err)
			}
		}
		// 更新成功
		return nil
	})
}

// CancelInterv 取消学生的面试预约
func (*Stu) CancelInterv(netid string, intervid int) error {
	// 查询面试记录
	var record model.Interv
	if err := model.DB.Model(&model.Interv{}).Where("id = ?", intervid).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("没找到"), common.NotFoundErr)
		}
		logger.GinLogger.Errorf("查询面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 检查是否在不可预约时间段内
	for timeRange, _ := range BlockTable {
		if record.Time.After(timeRange.Start) && record.Time.Before(timeRange.End) {
			return common.ErrNew(errors.New("当前面试信息无法修改"), common.AuthErr)
		}
	}
	if time.Now().After(record.Time) {
		return common.ErrNew(errors.New("面试时间已过"), common.OpErr)
	}
	// 检查面试记录是否属于该学生
	if record.NetID == nil && *record.NetID != netid {
		return common.ErrNew(errors.New("该面试记录不属于该学生"), common.AuthErr)
	}

	// 检查面试时间是否在半小时内或已经错过
	if record.Time.Before(time.Now().Add(30*time.Minute)) || record.Time.Before(time.Now()) {
		return common.ErrNew(errors.New("面试时间在半小时内或已经错过，无法取消预约"), common.AuthErr)
	}

	// 取消学生的面试预约
	if err := model.DB.Model(&model.Interv{}).Where("id = ?", intervid).Update("netid", nil).Error; err != nil {
		logger.GinLogger.Errorf("取消预约面试失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 更新成功
	return nil
}

// GetRes 获取学生的面试结果
func (*Stu) GetRes(netid string) (model.Interv, error) {
	if time.Now().Before(AvailableTime) {
		return model.Interv{}, common.ErrNew(errors.New("面试结果尚未公布"), common.OpErr)
	}
	var data model.Interv
	if err := model.DB.Model(&model.Interv{}).Where("netid = ?", netid).First(&data).Error; err != nil {
		logger.GinLogger.Errorf("查询学生面试结果失败: %v", err)
		return model.Interv{}, common.ErrNew(err, common.SysErr)
	}
	return data, nil
}
