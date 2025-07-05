package service

import (
	"errors"
	"zhaoxin2025/common"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"

	"gorm.io/gorm"
)

type Que struct{}

// 获取问题
// 支持通过问题内容(que)、部门(department)、链接(url)进行筛选，并进行分页(pager)
func (*Que) Get(que string, department []model.Department, url string) (int64, []model.Que, error) {
	var data []model.Que
	var count int64

	db := model.DB.Model(&model.Que{}) // 初始化查询构建器，针对 Que 模型
	// 如果提供了问题内容筛选条件，则添加模糊匹配
	if que != "" {
		db = db.Where("question LIKE ?", "%"+que+"%")
	}
	// 如果提供了部门筛选条件，则添加精确匹配
	if department != nil {
		db = db.Where("department IN ?", department)
	}
	if url != "" {
		db = db.Where("url LIKE ?", url)
	}
	if err := db.Count(&count).Error; err != nil {
		logger.DatabaseLogger.Errorf("统计问题数量失败：%v", err)
		return 0, nil, err
	}
	// 执行分页查询并获取结果
	if err := db.Omit("created_at", "updated_at", "deleted_at").Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取问题失败：%v", err)
		return 0, nil, common.ErrNew(err, common.SysErr)
	}
	return count, data, nil
}

// New 新建问题
// 批量创建问题记录
func (*Que) New(list []model.Que) error {
	if err := model.DB.Model(&model.Que{}).Create(&list).Error; err != nil {
		logger.DatabaseLogger.Errorf("创建问题失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// Delete 删除问题
// 根据提供的ID列表批量删除问题记录
func (*Que) Delete(ids []int) error {
	// 直接使用 GORM 的 Delete 方法，通过主键ID列表删除记录
	if err := model.DB.Model(&model.Que{}).Delete(&model.Que{}, ids).Error; err != nil {
		logger.DatabaseLogger.Errorf("批量删除问题失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// Update 更新问题
// 根据提供的ID和信息更新单个问题记录
func (*Que) Update(info model.Que) error {
	// 使用 GORM 的 Updates 方法更新记录
	var record model.Que
	if err := model.DB.Model(&model.Que{}).Where("id = ?", info.ID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("该问题不存在"), common.NotFoundErr)
		}
		logger.DatabaseLogger.Errorf("查询问题失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	if err := model.DB.Model(&model.Que{}).Where("id = ?", info.ID).Updates(&info).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新问题失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// LuckyDog 指定问题
// 根据提供的NetID和QueID来给幸运儿指定问题
func (*Que) LuckyDog(netid string, queid int) error {
	var record model.Stu
	if err := model.DB.Model(&model.Stu{}).Where("netid = ?", netid).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("禁止虚空索敌"), common.NotFoundErr)
		} else {
			logger.DatabaseLogger.Errorf("查询面试记录失败：%v", err)
			return common.ErrNew(err, common.SysErr)
		}
	}
	if err := model.DB.Model(&model.Stu{}).Where("netid = ?", netid).Update("queid", queid).Error; err != nil {
		logger.DatabaseLogger.Errorf("指定问题失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}
