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
func (*Interv) Get(info map[string]interface{}) ([]model.Interv, error) {
	var data []model.Interv
	db := model.DB.Model(&model.Interv{}).Where(info)

	// 根据日期筛选
	if info["date"] != nil {
		date := info["date"].(time.Time)
		dayRange := DayRange(date)
		db = db.Where("time BETWEEN ? AND ?", dayRange.Start, dayRange.End)
	}

	// 执行分页查询
	page := info["page"].(int)
	limit := info["limit"].(int)
	if err := db.Offset((page - 1) * limit).Limit(limit).Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}

	// 检查是否有查询结果
	if len(data) == 0 {
		return nil, common.ErrNew(errors.New("没有查询到面试记录"), common.OpErr)
	}
	return data, nil
}

// New 创建新的面试时间记录，返回冲突的时间列表
func (*Interv) New(info []time.Time) ([]time.Time, error) {
	var conflict []time.Time
	for _, v := range info {
		newInterv := model.Interv{
			Time: v,
		}
		var count int64
		// 检查时间是否已存在
		if err := model.DB.Model(&model.Interv{}).Where("time = ?", v).Count(&count).Error; err != nil {
			logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
			return nil, common.ErrNew(err, common.SysErr)
		}
		// 如果时间冲突，记录并跳过
		if count > 0 {
			conflict = append(conflict, v)
			continue
		}
		// 创建新的面试记录
		if err := model.DB.Model(&model.Interv{}).Create(&newInterv).Error; err != nil {
			logger.DatabaseLogger.Errorf("创建面试记录失败: %v", err)
			return nil, common.ErrNew(err, common.SysErr)
		}
	}
	return conflict, nil
}

// Update 更新面试记录信息
func (*Interv) Update(info map[string]interface{}) error {
	netid := info["netid"].(string)
	var count int64
	// 检查要更新的记录是否存在
	if err := model.DB.Model(&model.Interv{}).Where("net_id = ?", netid).Count(&count).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	if count == 0 {
		return common.ErrNew(errors.New("面试记录不存在"), common.OpErr)
	}
	// // 执行更新操作
	if err := model.DB.Model(&model.Interv{}).Where("net_id = ?", netid).Updates(info).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// Delete 批量删除面试记录，使用事务确保操作的原子性
func (*Interv) Delete(info []int) ([]int, error) {
	var fail []int
	// 使用事务进行批量删除
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range info {
			var count int64
			// 检查记录是否存在
			if err := tx.Model(&model.Interv{}).Where("id = ?", id).Count(&count).Error; err != nil {
				logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
				return common.ErrNew(err, common.SysErr)
			}
			// 记录不存在，加入失败列表但继续处理其他ID
			if count == 0 {
				fail = append(fail, id)
				continue
			}
			// 执行删除操作
			if err := tx.Delete(&model.Interv{}, id).Error; err != nil {
				logger.DatabaseLogger.Errorf("删除面试记录失败 (ID=%d): %v", id, err)
				return common.ErrNew(err, common.SysErr)
			}
		}
		return nil
	})
	// 返回删除失败的ID列表和可能的错误
	return fail, err
}
