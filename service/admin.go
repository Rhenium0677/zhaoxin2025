package service

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"zhaoxin2025/common"
	"zhaoxin2025/config"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"

	"gorm.io/gorm"
)

var AvailableTime = time.Date(2026, time.October, 1, 0, 0, 0, 0, time.Local)
var BlockTable = make(map[TimeRange]struct{})

type Admin struct{}

func (*Admin) Login(netid string, password string) (model.Admin, error) {
	// 检查管理员是否存在
	var record model.Admin
	if err := model.DB.Model(&model.Admin{}).Where("netid = ?", netid).First(&record).Error; err != nil {
		return model.Admin{}, common.ErrNew(err, common.OpErr)
	}
	// 加密传入的密码
	decrypted, err := Decrypt(record.Password)
	if err != nil {
		return model.Admin{}, common.ErrNew(err, common.SysErr)
	}
	// 检查密码是否正确
	if decrypted != password {
		return model.Admin{}, common.ErrNew(errors.New("用户名或密码错误"), common.OpErr)
	}
	return record, nil
}

func (*Admin) Update(netid string, name string, password string) error {
	info := model.Admin{
		Name:     name,
		Password: password,
	}
	// 检查管理员是否存在
	var count int64
	if err := model.DB.Model(&model.Admin{}).Where("netid = ?", netid).Count(&count).Error; err != nil {
		logger.DatabaseLogger.Errorf("检查管理员是否存在失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	if count == 0 {
		return common.ErrNew(errors.New("管理员不存在"), common.OpErr)
	}
	// 更新管理员信息
	if err := model.DB.Model(&model.Admin{}).Where("netid = ?", netid).Updates(&info).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新管理员信息失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// GetStu 筛选并获取学生信息
func (*Admin) GetStu(stuInfo model.Stu, intervInfo model.Interv, page int, limit int) ([]model.Stu, int64, error) {
	var (
		data  []model.Stu
		count int64
	)

	// 基础学生筛选
	db := model.DB.Model(&model.Stu{}).
		Where("stus.netid <> stus.openid").
		Where(&stuInfo)

	// 是否存在面试条件
	hasIntervCond := false

	// EXISTS 子查询：只在提供面试筛选条件时启用
	exists := model.DB.Model(&model.Interv{}).
		Select("1").
		Where("intervs.netid = stus.netid")

	if s := strings.TrimSpace(intervInfo.Evaluation); s != "" {
		exists = exists.Where("intervs.evaluation LIKE ?", "%"+s+"%")
		hasIntervCond = true
	}
	if s := strings.TrimSpace(intervInfo.Interviewer); s != "" {
		exists = exists.Where("intervs.interviewer LIKE ?", "%"+s+"%")
		hasIntervCond = true
	}
	if intervInfo.Star != 0 {
		exists = exists.Where("intervs.star = ?", intervInfo.Star)
		hasIntervCond = true
	}
	// 传入 1 表示通过，2 表示未通过；0 表示不筛选
	if intervInfo.Pass == 1 || intervInfo.Pass == 2 {
		passVal := 1
		if intervInfo.Pass == 2 {
			passVal = 0
		}
		exists = exists.Where("intervs.pass = ?", passVal)
		hasIntervCond = true
	}

	if hasIntervCond {
		db = db.Where("EXISTS (?)", exists)
	}

	// 统计总数（无需 DISTINCT）
	if err := db.Count(&count).Error; err != nil {
		logger.DatabaseLogger.Errorf("统计学生信息失败：%v", err)
		return nil, 0, common.ErrNew(err, common.SysErr)
	}

	offset := (page - 1) * limit

	// 查询数据；预加载面试记录
	if err := db.
		Preload("Interv").
		Order("stus.id ASC").
		Offset(offset).
		Limit(limit).
		Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取学生信息失败：%v", err)
		return nil, 0, common.ErrNew(err, common.SysErr)
	}

	return data, count, nil
}

// UpdateStu 更新一个学生信息
func (*Admin) UpdateStu(stuInfo model.Stu, intervInfo model.Interv) error {
	// 检查学生是否存在
	var existed model.Stu
	if err := model.DB.Model(&model.Stu{}).Where("id = ?", stuInfo.ID).First(&existed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("学生不存在"), common.NotFoundErr)
		}
		logger.DatabaseLogger.Errorf("检查学生是否存在失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return model.DB.Transaction(func(tx *gorm.DB) error {
		// 更新学生信息
		if err := tx.Model(&model.Stu{}).Where("id = ?", stuInfo.ID).Updates(&stuInfo).Error; err != nil {
			logger.DatabaseLogger.Errorf("更新学生信息失败：%v", err)
			return common.ErrNew(err, common.SysErr)
		}
		// 更新面试信息
		intervService := &Interv{}
		if err := intervService.Update(intervInfo); err != nil {
			logger.DatabaseLogger.Errorf("更新面试信息失败：%v", err)
			return err
		}
		return nil
	})
}

func (*Admin) DeleteStu(id int64) error {
	// 检查学生是否存在
	var existed model.Stu
	if err := model.DB.Model(&model.Stu{}).Where("id = ?", id).First(&existed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("学生不存在"), common.NotFoundErr)
		}
		logger.DatabaseLogger.Errorf("检查学生是否存在失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return model.DB.Transaction(func(tx *gorm.DB) error {
		// 删除学生信息
		if err := tx.Model(&model.Stu{}).Where("id = ?", id).Delete(&model.Stu{}).Error; err != nil {
			logger.DatabaseLogger.Errorf("删除学生信息失败：%v", err)
			return common.ErrNew(err, common.SysErr)
		}
		// 删除面试信息
		if err := tx.Model(&model.Interv{}).Where("netid = ?", existed.NetID).Delete(&model.Interv{}).Error; err != nil {
			logger.DatabaseLogger.Errorf("删除面试信息失败：%v", err)
			return common.ErrNew(err, common.SysErr)
		}
		return nil
	})
}

// Excelize 获取学生信息并导出为excel
func (*Admin) Excelize() error {
	// 获取所有学生信息
	var data []model.Stu
	if err := model.DB.Model(&model.Stu{}).Preload("Interv").Where("netid LIKE ?", "2%").Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取学生信息失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	// 创建excel文件
	if err := Excelize(data, "tenzor2025.xlsx"); err != nil {
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// Stat 获取学生信息并统计
func (*Admin) Stat() (Stat, error) {
	var stus []model.Stu
	if err := model.DB.Model(&model.Stu{}).Find(&stus).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取学生信息失败：%v", err)
		return Stat{}, common.ErrNew(err, common.SysErr)
	}
	var intervs []model.Interv
	if err := model.DB.Model(&model.Interv{}).Find(&intervs).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取面试信息失败：%v", err)
		return Stat{}, common.ErrNew(err, common.SysErr)
	}
	return GetStat(stus, intervs), nil
}

// Register 管理员注册
func (*Admin) Register(netid string, name string, password string, level model.AdminLevel) error {
	// 检查管理员是否存在
	var record model.Admin
	if err := model.DB.Model(&model.Admin{}).Where("netid = ?", netid).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用data结构体存储管理员信息
			// 加密密码
			encrypted, err := Encrypt(password)
			if err != nil {
				return common.ErrNew(err, common.SysErr)
			}
			data := model.Admin{
				NetID:    netid,
				Name:     name,
				Password: encrypted,
				Level:    level,
			}
			// 在数据库中插入数据
			if err := model.DB.Model(&model.Admin{}).Create(&data).Error; err != nil {
				logger.DatabaseLogger.Errorf("插入管理员信息失败：%v", err)
				return common.ErrNew(err, common.SysErr)
			}
			return nil

		}
		logger.DatabaseLogger.Errorf("检查管理员是否存在失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return common.ErrNew(errors.New("管理员已存在"), common.OpErr)
}

func (*Admin) SendResultMessage() error {
	// 获取所有学生信息
	var data []model.Stu
	if err := model.DB.Model(&model.Stu{}).Where("message > 3").Preload("Interv").Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取学生信息失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	for _, stu := range data {
		if stu.Message < 4 || stu.Interv == nil {
			continue
		}
		err := SendResult(stu)
		if err != nil {
			println("添加面试结果订阅消息失败: %v，学生netid: %s", err, stu.NetID)
		}
	}
	return nil
}

// 发送短信
// 有关短信的操作
func (a *Admin) AliyunSendItvResMsg() (cocacola interface{}, err error) {
	// 向Aliyun发送请求
	var user []model.Stu
	if err = model.DB.
		Model(&model.Stu{}).Preload("Interv").
		Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNew(errors.New("当前学号不存在"), common.NotFoundErr)
		}
		return nil, common.ErrNew(errors.New("没有查到,请联系管理员解决"), common.SysErr)
	}
	counter := 0
	for _, value := range user {
		if value.Interv == nil {
			continue
		}
		if counter%10 == 0 && counter != 0 {
			time.Sleep(3 * time.Second) // 每发送10条短信，休息3秒
		}
		var check bool
		if value.Interv.Pass == 1 {
			check = true
		} else if value.Interv.Pass == 0 {
			check = false
		} else {
			continue
		}
		// 发短信在这里
		err := sendItvResMsg(check, value.Phone, value.Name, DepartToChinese(value.Depart))
		if err != nil {
			logger.DatabaseLogger.Errorf("name: %s 发送短信失败: %v", err)
			continue
		}
		counter += 1
		logger.DatabaseLogger.Infof("name: %s 发送短信成功: %s", value.Phone)
	}
	return nil, nil
}

// 这是发送面试结果短信的函数
func sendItvResMsg(pass bool, number string, name string, department string) error {
	client, _err := CreateClient(tea.String(config.Config.AlibabaCloudAccessKeyID), tea.String(config.Config.AlibabaCloudAccessKeySecret))
	if _err != nil {
		return _err
	}

	var tpCode string
	// 短信对应代码,结合自己的短信格式的代码进行设置

	if pass {
		tpCode = "SMS_495000002"
	} else {
		tpCode = "SMS_494950005"
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(number),
		SignName:      tea.String("西咸新区挑战信息科技"),
		TemplateCode:  tea.String(tpCode),
		TemplateParam: tea.String(fmt.Sprintf("{\"name\":\"%s\",\"department\":\"%s\"}", name, department)),
	}

	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, _err := client.SendSmsWithOptions(sendSmsRequest, runtime)
		//if jsonData, err := json.MarshalIndent(result, "", "  "); err == nil {
		//	println(string(jsonData))
		//} else {
		//	println(err)
		//}
		if _err != nil {
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		msg, _err := util.AssertAsString(error.Message)
		if _err != nil {
			return _err
		}
		return errors.New(*msg)
	}
	return _err
}

func AliyunSendItvTimeMsg() (fails []FailSend, err error) {
	now := time.Now()
	start := now.Add(20 * time.Minute)
	end := now.Add(30 * time.Minute)

	// 一次性联表查询，避免循环里再查
	var rows []struct {
		NetID  string           `gorm:"column:netid"`
		Time   time.Time        `gorm:"column:time"`
		Phone  string           `gorm:"column:phone"`
		Name   string           `gorm:"column:name"`
		Depart model.Department `gorm:"column:depart"`
	}
	if e := model.DB.
		Table("intervs").
		Select("stus.netid, intervs.time, stus.phone, stus.name, stus.depart").
		Joins("JOIN stus ON stus.netid = intervs.netid").
		Where("intervs.time BETWEEN ? AND ?", start, end).
		Find(&rows).Error; e != nil {
		return fails, common.ErrNew(errors.New("查询面试记录失败"), common.SysErr)
	}

	for _, r := range rows {
		sendErr := sendItvTimeMsg(r.Phone, r.Name, r.Time.Format("2006年1月2日 15:04"), DepartToChinese(r.Depart))
		if sendErr != nil {
			fails = append(fails, FailSend{NetID: r.NetID, ErrCode: -1})
			logger.DatabaseLogger.Errorf("发送面试时间短信失败 netid=%s err=%v", r.NetID, sendErr)
		}
	}

	return fails, nil
}

// 发送面试通知短信
func sendItvTimeMsg(phone string, name string, time string, department string) (err error) {
	client, err := CreateClient(&config.Config.AlibabaCloudAccessKeyID, &config.Config.AlibabaCloudAccessKeySecret)
	if err != nil {
		return err
	}
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String("西咸新区挑战信息科技"),
		TemplateCode:  tea.String("SMS_494950004"),
		TemplateParam: tea.String(fmt.Sprintf("{\"name\":\"%s\",\"time\":\"%s\",\"department\":\"%s\"}", name, time, department)),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, _err := client.SendSmsWithOptions(sendSmsRequest, runtime)
		//if jsonData, err := json.MarshalIndent(result, "", "  "); err == nil {
		//	println(string(jsonData))
		//} else {
		//	println(err)
		//}
		if _err != nil {
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		var err = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			err = _t
		} else {
			err.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 err
		msg, _err := util.AssertAsString(err.Message)
		if _err != nil {
			return _err
		}
		return errors.New(*msg)
	}
	return err
}

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}
