package service

import (
	"zhaoxin2025/common"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"

	"gorm.io/gorm"
)

type Que struct{}

// 获取问题
func (*Que) Get(que string, department model.Department, url string, pager common.PagerForm) ([]model.Que, error) {
	var data []model.Que
	db := model.DB.Model(&model.Que{})
	if que != "" {
		db = db.Where("question LIKE ?", "%"+que+"%")
	}
	if department != "" {
		db = db.Where("department = ?", department)
	}
	if url != "" {
		db = db.Where("url LIKE ?", "%"+url+"%")
	}
	if err := db.Offset((pager.Page - 1) * pager.Limit).Limit(pager.Limit).Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取问题失败：%v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	return data, nil
}

// 新建问题
func (*Que) New(list []model.Que) error {
	if err := model.DB.Model(&model.Que{}).Create(&list).Error; err != nil {
		logger.DatabaseLogger.Errorf("创建问题失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// 删除问题
func (*Que) Delete(ids []int) error {
	// 使用事务确保删除的原子性
	return model.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			if err := tx.Model(&model.Que{}).Where("id = ?", id).Delete(&model.Que{}).Error; err != nil {
				logger.DatabaseLogger.Errorf("删除问题失败，事务回滚：%v", err)
				return common.ErrNew(err, common.SysErr)
			}
		}
		return nil
	})
}

// 更新问题
func (*Que) Update(info UpdateQue) error {
	if err := model.DB.Model(&model.Que{}).Where("id = ?", info.ID).Updates(&info).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新问题失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}
