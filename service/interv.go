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

func (*Interv) New(info []time.Time) ([]time.Time, error) {
	var conflict []time.Time
	for _, v := range info {
		newInterv := model.Interv{
			Time: v,
		}
		var count int64
		if err := model.DB.Model(&model.Interv{}).Where("time = ?", v).Count(&count).Error; err != nil {
			logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
			return nil, common.ErrNew(err, common.SysErr)
		}
		if count > 0 {
			conflict = append(conflict, v)
			continue
		}
		if err := model.DB.Model(&model.Interv{}).Create(&newInterv).Error; err != nil {
			logger.DatabaseLogger.Errorf("创建面试记录失败: %v", err)
			return nil, common.ErrNew(err, common.SysErr)
		}
	}
	return conflict, nil
}

func (*Interv) Update(info IntervUpdate) error {
	var count int64
	if err := model.DB.Model(&model.Interv{}).Where("id = ?", info.ID).Count(&count).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	if count == 0 {
		return common.ErrNew(errors.New("面试记录不存在"), common.OpErr)
	}
	if err := model.DB.Model(&model.Interv{}).Where("id = ?", info.ID).Updates(&info).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新面试记录失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

func (*Interv) Get(info SelectStu) ([]model.Interv, error) {
	var data []model.Interv
	db := model.DB.Model(&model.Interv{})
	if info.NetID != "" {
		db = db.Where("netid = ?", info.NetID)
	}
	if info.Name != "" {
		db = db.Where("name = ?", info.Name)
	}
	if info.Phone != "" {
		db = db.Where("phone = ?", info.Phone)
	}
	if info.School != "" {
		db = db.Where("school = ?", info.School)
	}
	if info.First != "" {
		db = db.Where("first = ?", info.First)
	}
	if info.Second != "" {
		db = db.Where("second = ?", info.Second)
	}
	if info.Interviewer != "" {
		db = db.Where("interviewer = ?", info.Interviewer)
	}
	if info.Star != 0 {
		db = db.Where("star = ?", info.Star)
	}
	if err := db.Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	if len(data) == 0 {
		return nil, common.ErrNew(errors.New("没有查询到面试记录"), common.OpErr)
	}
	return data, nil
}

func (*Interv) Delete(info []int) ([]int, error) {
	var fail []int
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range info {
			var count int64
			if err := tx.Model(&model.Interv{}).Where("id = ?", id).Count(&count).Error; err != nil {
				logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
				return common.ErrNew(err, common.SysErr)
			}
			if count == 0 {
				fail = append(fail, id)
				continue
			}
			if err := tx.Delete(&model.Interv{}, id).Error; err != nil {
				logger.DatabaseLogger.Errorf("删除面试记录失败 (ID=%d): %v", id, err)
				return common.ErrNew(err, common.SysErr)
			}
		}
		return nil
	})
	return fail, err
}
