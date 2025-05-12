package service

import (
	"errors"
	"time"
	"zhaoxin2025/common"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"

	"gorm.io/gorm"
)

type Interv struct{}

// Get 按条件查询面试记录，支持分页
func (*Interv) Get(info map[string]any) ([]model.Interv, error) {
	var data []model.Interv
	db := model.DB.Model(&model.Interv{}).Where(info)

	// 根据日期筛选
	if info["date"] != nil {
		date := info["date"].(time.Time)
		dayRange := DayRange(date)
		db = db.Where("time BETWEEN ? AND ?", dayRange.Start, dayRange.End)
	}

	// 执行分页查询
	page, ok := info["page"].(int)
	if !ok {
		return nil, common.ErrNew(errors.New("分页参数错误"), common.ParamErr)
	}
	limit, ok := info["limit"].(int)
	if !ok {
		return nil, common.ErrNew(errors.New("分页参数错误"), common.ParamErr)
	}

	if err := db.Offset((page - 1) * limit).Limit(limit).Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	return data, nil
}

// New 创建新的面试时间记录
func (*Interv) New(info TimeRange, interval int) ([]time.Time, error) {
	var conflict []model.Interv
	if err := model.DB.Where("time BETWEEN ? AND ?", info.Start, info.End).Find(&conflict).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询冲突面试记录失败: %v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	if conflict != nil { // 如果有冲突，返回冲突的时间
		var conflictTimes []time.Time
		for _, record := range conflict {
			conflictTimes = append(conflictTimes, record.Time)
		}
		// 返回冲突的时间
		return conflictTimes, common.ErrNew(errors.New("冲突面试记录"), common.ConflictErr)
	}

	// 没有冲突，创建新的面试时间
	var intervs []model.Interv
	for t := info.Start; t.Before(info.End); t = t.Add(time.Duration(interval) * time.Minute) {
		intervs = append(intervs, model.Interv{
			Time: t,
		})
	}
	if err := model.DB.Create(&intervs).Error; err != nil {
		logger.DatabaseLogger.Errorf("创建面试时间失败: %v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	return nil, nil
}

// Update 更新面试记录信息
func (*Interv) Update(info model.Interv) error {
	var existed model.Interv
	// 检查要更新的记录是否存在
	if err := model.DB.Model(&model.Interv{}).Where("netid = ?", info.ID).First(&existed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := model.DB.Model(&model.Interv{}).Where("netid = ?", info.ID).Updates(&info).Error; err != nil {
				logger.DatabaseLogger.Errorf("更新面试记录失败: %v", err)
				return common.ErrNew(err, common.SysErr)
			}
			return nil
		}
		logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return common.ErrNew(errors.New("记录不存在"), common.NotFoundErr)
}

// Delete 批量删除面试记录
func (*Interv) Delete(info []int) error {
	if err := model.DB.Where("id in ?", info).Delete(&model.Interv{}).Error; err != nil {
		logger.DatabaseLogger.Errorf("删除面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 返回删除失败的ID列表和可能的错误
	return nil
}

func (*Interv) BlockAndRecover(timeRange TimeRange, block bool) error {
	if block {
		if err := model.DB.Where("time BETWEEN ? AND ?", timeRange.Start, timeRange.End).Delete(&model.Interv{}).Error; err != nil {
			logger.DatabaseLogger.Errorf("禁止面试失败: %v", err)
			return common.ErrNew(err, common.SysErr)
		}
	} else {
		if err := model.DB.Where("time BETWEEN ? AND ?", timeRange.Start, timeRange.End).Update("deleted_at", nil).Error; err != nil {
			logger.DatabaseLogger.Errorf("恢复面试失败: %v", err)
			return common.ErrNew(err, common.SysErr)
		}
	}
	return nil
}
