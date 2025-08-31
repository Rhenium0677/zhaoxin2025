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
func (*Interv) Get(info model.Interv, date time.Time, page int, limit int) (int64, []IntervInfo, error) {
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
	if limit == 0 {
		page = 1           // 如果 limit 为 0，则页码无意义，重置为 1
		limit = int(count) // 如果 limit 为 0，则查询所有记录
	}
	if err := db.Offset((page - 1) * limit).Limit(limit).Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
		return 0, nil, common.ErrNew(err, common.SysErr)
	}
	var dataInfo []IntervInfo
	for i, _ := range data {
		if data[i].NetID != nil {
			var stu model.Stu
			if err := model.DB.Model(&model.Stu{}).Where("netid = ?", *data[i].NetID).First(&stu).Error; err != nil {
				logger.DatabaseLogger.Errorf("查询学生信息失败: %v", err)
				return 0, nil, common.ErrNew(err, common.SysErr)
			} else {
				if stu.QueID != 0 && data[i].QueID != stu.QueID { // 如果学生已经抽过题且面试记录中的题目ID不匹配，则更新面试记录
					logger.DatabaseLogger.Warnf("面试记录中的题目ID与学生信息不匹配，正在修正")
					if err := model.DB.Model(&model.Interv{}).Where("id = ?", data[i].ID).Update("queid", stu.QueID).Error; err != nil {
						logger.DatabaseLogger.Errorf("更新面试问题ID失败: %v", err)
						return 0, nil, common.ErrNew(err, common.SysErr)
					}
					data[i].QueID = stu.QueID // 关联学生信息
				}
			}
			dataInfo = append(dataInfo, IntervInfo{
				Interv: data[i],
				Name:   stu.Name,
			})
			continue
		}
		dataInfo = append(dataInfo, IntervInfo{
			Interv: data[i],
			Name:   "",
		})
	}
	return count, dataInfo, nil
}

// GetDate 获取所有日期对应的面试个数
func (*Interv) GetDate() (map[string]int, error) {
	type row struct {
		D time.Time `gorm:"column:d"`
		C int       `gorm:"column:c"`
	}
	var rows []row
	if err := model.DB.Model(&model.Interv{}).Where("netid IS NOT NULL").
		Select("DATE(`time`) AS d, COUNT(*) AS c").
		Group("DATE(`time`)").
		Order("d").Scan(&rows).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询可用面试日期失败: %v", err)
		return nil, common.ErrNew(err, common.SysErr)
	}
	intervDate := make(map[string]int, len(rows))
	for _, r := range rows {
		intervDate[Date(r.D)] = r.C
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

func (*Interv) Cancel(id int) error {
	tx := model.DB.Begin()
	if err := tx.Model(&model.Stu{}).Where("netid = (SELECT netid FROM intervs WHERE id = ?)", id).
		Update("queid", 0).Error; err != nil {
		logger.DatabaseLogger.Errorf("重置学生问题ID失败: %v", err)
		tx.Rollback()
		return common.ErrNew(err, common.SysErr)
	}
	if err := tx.Model(&model.Interv{}).Where("id = ?", id).
		Updates(map[string]interface{}{"netid": nil, "status": 0, "quetime": 0, "queid": 0, "evaluation": "", "star": 0, "pass": 0, "interviewer": "", "department": "none"}).Error; err != nil {
		logger.DatabaseLogger.Errorf("取消面试记录失败: %v", err)
		tx.Rollback()
		return common.ErrNew(err, common.SysErr)
	}
	if err := tx.Commit().Error; err != nil {
		logger.DatabaseLogger.Errorf("提交事务失败: %v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// Update 更新面试记录信息
func (*Interv) Update(info model.Interv) error {
	var record model.Interv
	// 检查要更新的记录是否存在
	if info.ID != 0 {
		if err := model.DB.Model(&model.Interv{}).Where("id = ?", info.ID).First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return common.ErrNew(errors.New("禁止虚空索敌"), common.NotFoundErr)
			}
			logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
			return common.ErrNew(err, common.SysErr)
		}
	} else {
		if err := model.DB.Model(&model.Interv{}).Where("netid = ?", info.NetID).First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return common.ErrNew(errors.New("禁止虚空索敌"), common.NotFoundErr)
			}
			logger.DatabaseLogger.Errorf("查询面试记录失败: %v", err)
			return common.ErrNew(err, common.SysErr)
		}
		info.ID = record.ID
	}
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if info.Evaluation != "" {
			if err := tx.Model(&model.Interv{}).Where("id = ?", info.ID).Update("status", 2).Error; err != nil {
				logger.DatabaseLogger.Errorf("更新面试状态失败: %v", err)
				return common.ErrNew(err, common.SysErr)
			}
		}
		if info.Pass == 2 {
			if err := tx.Model(&model.Interv{}).Where("id = ?", info.ID).Update("pass", 0).Error; err != nil {
				logger.DatabaseLogger.Errorf("更新面试通过状态失败: %v", err)
				return common.ErrNew(err, common.SysErr)
			}
		}
		if err := tx.Model(&model.Interv{}).Where("id = ?", info.ID).Updates(&info).Error; err != nil {
			logger.DatabaseLogger.Errorf("更新面试记录失败: %v", err)
			return common.ErrNew(err, common.SysErr)
		}
		return nil
	})
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
	var student model.Stu
	// 查询学生信息
	if err := model.DB.Model(&model.Stu{}).Where("netid = ?", netid).First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Que{}, common.ErrNew(errors.New("没有找到学生信息"), common.NotFoundErr)
		}
		logger.DatabaseLogger.Errorf("查询学生信息失败: %v", err)
		return model.Que{}, common.ErrNew(err, common.SysErr)
	}

	// 检查学生是否已经抽过题目且题目仍然有效
	if student.QueID != 0 {
		if que, err := i.validateExistingQuestion(student.QueID, department); err == nil {
			// 题目有效，更新相关记录并返回
			return que, i.updateRecordsForExistingQuestion(netid, student.QueID, timeRecord)
		}
		// 题目无效或部门不匹配，需要重新抽题
	}

	// 抽取新题目
	return i.drawNewQuestion(netid, department, timeRecord)
}

// validateExistingQuestion 验证已存在的题目是否有效且部门匹配
func (*Interv) validateExistingQuestion(queID int, department model.Department) (model.Que, error) {
	var que model.Que
	if err := model.DB.Model(&model.Que{}).Where("id = ?", queID).First(&que).Error; err != nil {
		return model.Que{}, err
	}
	if que.Department != department {
		return model.Que{}, errors.New("部门不匹配")
	}
	return que, nil
}

// updateRecordsForExistingQuestion 为已存在的题目更新相关记录
func (*Interv) updateRecordsForExistingQuestion(netid string, queID int, timeRecord int64) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		// 更新面试状态
		if err := tx.Model(&model.Interv{}).Where("netid = ?", netid).
			Updates(map[string]interface{}{
				"status":  1,
				"quetime": timeRecord,
				"queid":   queID,
			}).Error; err != nil {
			logger.DatabaseLogger.Errorf("更新面试状态失败: %v", err)
			return common.ErrNew(err, common.SysErr)
		}

		// 增加题目被抽中次数
		if err := tx.Model(&model.Que{}).Where("id = ?", queID).
			Update("times", gorm.Expr("times + ?", 1)).Error; err != nil {
			logger.DatabaseLogger.Errorf("更新问题被抽中次数失败: %v", err)
			return common.ErrNew(err, common.SysErr)
		}

		return nil
	})
}

// drawNewQuestion 抽取新题目
func (*Interv) drawNewQuestion(netid string, department model.Department, timeRecord int64) (model.Que, error) {
	// 查询该部门的所有题目
	var questions []model.Que
	if err := model.DB.Model(&model.Que{}).Where("department = ?", department).Find(&questions).Error; err != nil {
		logger.DatabaseLogger.Errorf("查询问题失败: %v", err)
		return model.Que{}, common.ErrNew(err, common.SysErr)
	}

	if len(questions) == 0 {
		return model.Que{}, common.ErrNew(errors.New("没有找到问题"), common.NotFoundErr)
	}

	// 随机选择题目
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	selectedQuestion := questions[r.Intn(len(questions))]

	// 在事务中更新所有相关记录
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 更新学生的题目ID
		if err := tx.Model(&model.Stu{}).Where("netid = ?", netid).
			Update("queid", selectedQuestion.ID).Error; err != nil {
			logger.DatabaseLogger.Errorf("更新学生问题ID失败: %v", err)
			return common.ErrNew(err, common.SysErr)
		}

		// 更新面试状态
		if err := tx.Model(&model.Interv{}).Where("netid = ?", netid).
			Updates(map[string]interface{}{
				"status":  1,
				"quetime": timeRecord,
				"queid":   selectedQuestion.ID,
			}).Error; err != nil {
			logger.DatabaseLogger.Errorf("更新面试状态失败: %v", err)
			return common.ErrNew(err, common.SysErr)
		}

		// 增加题目被抽中次数
		if err := tx.Model(&model.Que{}).Where("id = ?", selectedQuestion.ID).
			Update("times", gorm.Expr("times + ?", 1)).Error; err != nil {
			logger.DatabaseLogger.Errorf("更新问题被抽中次数失败: %v", err)
			return common.ErrNew(err, common.SysErr)
		}

		return nil
	})

	if err != nil {
		return model.Que{}, err
	}

	return selectedQuestion, nil
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
