package service

import (
	"zhaoxin2025/common"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"
)

type Que struct{}

// 获取问题
// 支持通过问题内容(que)、部门(department)、链接(url)进行筛选，并进行分页(pager)
func (*Que) Get(que string, department model.Department, url string, pager common.PagerForm) ([]model.Que, error) {
	var data []model.Que
	db := model.DB.Model(&model.Que{}) // 初始化查询构建器，针对 Que 模型
	// 如果提供了问题内容筛选条件，则添加模糊匹配
	if que != "" {
		db = db.Where("question LIKE ?", "%"+que+"%")
	}
	// 如果提供了部门筛选条件，则添加精确匹配
	if department != "" {
		db = db.Where("department = ?", department)
	}
	// 如果提供了URL筛选条件，则添加模糊匹配
	if url != "" {
		db = db.Where("url LIKE ?", "%"+url+"%")
	}
	// 执行分页查询并获取结果
	if err := db.Omit("created_at", "updated_at", "deleted_at").Offset((pager.Page - 1) * pager.Limit).Limit(pager.Limit).Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取问题失败：%v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	return data, nil
}

// 新建问题
// 批量创建问题记录
func (*Que) New(list []model.Que) error {
	if err := model.DB.Model(&model.Que{}).Create(&list).Error; err != nil {
		logger.DatabaseLogger.Errorf("创建问题失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// 删除问题
// 根据提供的ID列表批量删除问题记录
func (*Que) Delete(ids []int) error {
	// 直接使用 GORM 的 Delete 方法，通过主键ID列表删除记录
	result := model.DB.Delete(&model.Que{}, ids)
	if result.Error != nil {
		logger.DatabaseLogger.Errorf("批量删除问题失败：%v", result.Error)
		return common.ErrNew(result.Error, common.SysErr)
	}
	return nil
}

// 更新问题
// 根据提供的ID和信息更新单个问题记录
func (*Que) Update(info map[string]interface{}) error {
	// 从传入的map中获取问题ID
	id := info["id"]
	// 使用 GORM 的 Updates 方法，根据ID更新问题信息
	// Updates 方法会忽略 map 中模型不存在的字段
	if err := model.DB.Model(&model.Que{}).Where("id = ?", id).Updates(&info).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新问题失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}
