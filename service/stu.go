package service

import (
	"zhaoxin2025/common"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"
)

type Stu struct{}

func (*Stu) Login(netid string) error {
	//
	//
	return nil
}

func (*Stu) Update(info StuUpdate) error {
	//
	//
	if err := model.DB.Model(&model.Stu{}).Where("netid = ?", info.NetID).Updates(info).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新学生信息失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}
