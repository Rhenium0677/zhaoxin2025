package service

import (
	"zhaoxin2025/common"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"
)

type Que struct{}

func (*Que) Get() error {
	return nil
}

func (*Que) New(list []model.Que) error {
	if err := model.DB.Model(&model.Que{}).Create(&list).Error; err != nil {
		logger.DatabaseLogger.Errorf("创建问题失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

func (*Que) Delete(id int) error {
	if err := model.DB.Model(&model.Que{}).Where("id = ?", id).Delete(&model.Que{}).Error; err != nil {
		logger.DatabaseLogger.Errorf("删除问题失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

func (*Que) Update(id int, que model.Que) error {
	if err := model.DB.Model(&model.Que{}).Where("id = ?", id).Updates(&que).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新问题失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}
