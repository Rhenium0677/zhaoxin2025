package service

import (
	"errors"
	"math/rand"
	"time"
	"zhaoxin2025/common"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"

	"gorm.io/gorm"
)

type Interv struct{}

// Get 按条件查询面试记录，支持分页
func (*Interv) Get(info model.Interv, date time.Time, page int, limit int) (int64, []model.Interv, error) {
	var data []model.Interv
	var count int64
	db := model.DB.Model(&model.Interv{}).Where(&info)
	if date != (time.Time{}) {
		timeRange := DayRange(date)
		db = db.Where("time BETWEEN ? AND ?", timeRange.Start, timeRange.End)
	}
	if err := db.Count(&count).Error; err != nil {
		logger.DatabaseLogger.Errorf("统计数据总数失败: %v", err)
		return 0, nil, common.ErrNew(err, common.SysErr)
	}
	if err := db.Offset((page - 1) * limit).Limit(limit).Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
		return 0, nil, common.ErrNew(err, common.SysErr)
	}
	return count, data, nil
}

// GetDate 获取所有日期对应的面试个数
func (*Interv) GetDate() (map[string]int, error) {
	var data []model.Interv
	if err := model.DB.Model(&model.Interv{}).Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询可用面试日期失败: %v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	intervDate := make(map[string]int)
	for _, interv := range data {
		intervDate[Date(interv.Time)]++
	}
	return intervDate, nil
}

// New 创建新的面试时间记录
func (*Interv) New(info TimeRange, interval int) ([]time.Time, error) {
	var conflict []time.Time
	if err := model.DB.Model(&model.Interv{}).Where("time BETWEEN ? AND ?", info.Start, info.End).Pluck("time", &conflict).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询冲突面试记录失败: %v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	if len(conflict) > 0 { // 如果有冲突，返回冲突的时间
		// 返回冲突的时间
		return conflict, common.ErrNew(errors.New("有冲突面试记录"), common.ConflictErr)
	}

	// 没有冲突，创建新的面试时间
	var intervs []model.Interv
	for t := info.Start; t.Before(info.End); t = t.Add(time.Duration(interval) * time.Minute) {
		intervs = append(intervs, model.Interv{
			Time: t,
		})
	}
	if err := model.DB.Model(&model.Interv{}).Create(&intervs).Error; err != nil {
		logger.DatabaseLogger.Errorf("创建面试时间失败: %v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	return nil, nil
}

func (*Interv) Create(netid string, department model.Department, time time.Time) error {
	var interv = model.Interv{
		NetID:      &netid,
		Department: department,
		Time:       time,
	}
	if err := model.DB.Model(&model.Interv{}).Create(&interv).Error; err != nil {
		logger.DatabaseLogger.Errorf("创建面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// Update 更新面试记录信息
func (*Interv) Update(info model.Interv) error {
	var record model.Interv
	// 检查要更新的记录是否存在
	if err := model.DB.Model(&model.Interv{}).Where("id = ?", info.ID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("禁止虚空索敌"), common.NotFoundErr)
		}
		logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	if err := model.DB.Model(&model.Interv{}).Where("id = ?", info.ID).Updates(&info).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	if info.Evaluation != "" {
		if err := model.DB.Model(&model.Interv{}).Where("id = ?", info.ID).Update("status", 2).Error; err != nil {
			logger.DatabaseLogger.Errorf("更新面试状态失败: %v", err)
			return common.ErrNew(err, common.SysErr)
		}
	}
	return nil
}

// Delete 批量删除面试记录
func (*Interv) Delete(info []int) error {
	var count int64
	if err := model.DB.Model(&model.Interv{}).Where("id IN (?)", info).Count(&count).Error; err != nil {
		logger.DatabaseLogger.Errorf("统计要删除的面试记录数量失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	if count < int64(len(info)) {
		return common.ErrNew(errors.New("部分面试记录不存在"), common.NotFoundErr)
	}
	if err := model.DB.Model(&model.Interv{}).Where("id in ?", info).Delete(&model.Interv{}).Error; err != nil {
		logger.DatabaseLogger.Errorf("删除面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 返回删除失败的ID列表和可能的错误
	return nil
}

// Swap 交换两个面试记录
func (*Interv) Swap(id1, id2 int) error {
	var record1, record2 model.Interv
	// 检查要交换的记录是否存在
	if err := model.DB.Model(&model.Interv{}).Where("id = ?", id1).First(&record1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("禁止虚空索敌"), common.NotFoundErr)
		}
		logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	if err := model.DB.Model(&model.Interv{}).Where("id = ?", id2).First(&record2).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("禁止虚空索敌"), common.NotFoundErr)
		}
		logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	record1.ID, record2.ID = record2.ID, record1.ID // 交换ID
	return model.DB.Model(&model.Interv{}).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id1).Updates(&record1).Error; err != nil {
			logger.DatabaseLogger.Errorf("更新面试记录失败: %v", err)
			return common.ErrNew(err, common.SysErr)
		}
		if err := tx.Where("id = ?", id2).Updates(&record2).Error; err != nil {
			logger.DatabaseLogger.Errorf("更新面试记录失败: %v", err)
			return common.ErrNew(err, common.SysErr)
		}
		return nil
	})
}

// GetQue 为一个学生抽题，若幸运儿则假装抽过
func (i *Interv) GetQue(netid string, department model.Department, timeRecord int64) (model.Que, error) {
	var record model.Stu
	if err := model.DB.Model(&model.Stu{}).Where("netid = ?", netid).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Que{}, common.ErrNew(errors.New("没有找到学生信息"), common.NotFoundErr)
		}
		logger.DatabaseLogger.Errorf("查询学生信息失败: %v", err)
		return model.Que{}, common.ErrNew(err, common.SysErr)
	}
	reshuffle := false
	if record.QueID != 0 {
		var que model.Que
		if err := model.DB.Model(&model.Que{}).Where("queid = ?", record.QueID).First(&que).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				reshuffle = true
			} else {
				logger.DatabaseLogger.Errorf("查询问题失败: %v", err)
				return model.Que{}, common.ErrNew(err, common.SysErr)
			}
		}
		if que.Department == department && !reshuffle {
			return que, nil // 幸运儿，若部门匹配则直接返回，否则重抽
		}
	}
	var data []model.Que
	if err := model.DB.Model(&model.Que{}).Where("department = ?", department).Find(&data).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Que{}, common.ErrNew(errors.New("没有找到问题"), common.NotFoundErr)
		}
		logger.DatabaseLogger.Errorf("查询问题失败: %v", err)
		return model.Que{}, common.ErrNew(err, common.SysErr)
	}
	if len(data) == 0 {
		return model.Que{}, common.ErrNew(errors.New("没有找到问题"), common.NotFoundErr)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	queid := r.Intn(len(data))
	tx := model.DB.Begin()
	// 先更新 QueID
	if err := tx.Model(&model.Stu{}).Where("netid = ?", netid).
		Update("queid", data[queid].ID).Error; err != nil {
		tx.Rollback()
		logger.DatabaseLogger.Errorf("更新学生问题ID失败: %v", err)
		return model.Que{}, common.ErrNew(err, common.SysErr)
	}
	if err := tx.Model(&model.Interv{}).Where("netid = ?", netid).Updates(map[string]interface{}{"status": 1, "quetime": timeRecord, "queid": queid}).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新面试状态失败: %v", err)
		tx.Rollback()
		return model.Que{}, common.ErrNew(err, common.SysErr)
	}
	// 增加被抽中次数
	if err := tx.Model(&model.Que{}).Where("queid = ?", queid).
		Update("times", gorm.Expr("times + ?", 1)).
		Error; err != nil {
		logger.DatabaseLogger.Errorf("更新问题被抽中次数失败: %v", err)
		tx.Rollback()
		return model.Que{}, common.ErrNew(err, common.SysErr)
	}
	if err := tx.Commit().Error; err != nil {
		logger.DatabaseLogger.Errorf("提交事务失败: %v", err)
		return model.Que{}, common.ErrNew(err, common.SysErr)
	}
	return data[queid], nil
}

func (*Interv) BlockAndRecover(timeRange TimeRange, block bool) error {
	if block {
		BlockTable[timeRange] = struct{}{}
	} else {
		if _, ok := BlockTable[timeRange]; !ok {
			return common.ErrNew(errors.New("没有找到对应的封禁时间段"), common.NotFoundErr)
		}
		delete(BlockTable, timeRange)
	}
	return nil
}
