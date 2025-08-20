package service

import (
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"time"
	"zhaoxin2025/common"
	"zhaoxin2025/config"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"

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

// 筛选并获取学生信息
func (*Admin) GetStu(stuInfo model.Stu, intervInfo model.Interv, page int, limit int) ([]model.Stu, int64, error) {
	var data []model.Stu
	var count int64
	db := model.DB.Model(&model.Stu{}).Preload("Interv").Where(&stuInfo).Where("netid != openid")
	if intervInfo.Evaluation != "" {
		// 关联查询
		db = db.Where("interv.evaluation LIKE ?", intervInfo.Evaluation)
	}
	if intervInfo.Interviewer != "" {
		// 关联查询
		db = db.Where("interv.interviewer LIKE ?", intervInfo.Interviewer)
	}
	if err := db.Count(&count).Error; err != nil {
		logger.DatabaseLogger.Errorf("统计学生信息失败：%v", err)
		return nil, 0, common.ErrNew(err, common.SysErr)
	}
	if err := db.Offset((page - 1) * limit).Limit(limit).Find(&data).Error; err != nil {
		logger.DatabaseLogger.Errorf("获取学生信息失败：%v", err)
		return nil, 0, common.ErrNew(err, common.SysErr)
	}
	return data, count, nil
}

// 更新一个学生信息
func (*Admin) UpdateStu(stuInfo model.Stu) error {
	// 检查学生是否存在
	var existed model.Stu
	if err := model.DB.Model(&model.Stu{}).Where("id = ?", stuInfo.ID).First(&existed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNew(errors.New("学生不存在"), common.NotFoundErr)
		}
		logger.DatabaseLogger.Errorf("检查学生是否存在失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	if err := model.DB.Model(&model.Stu{}).Where("id = ?", stuInfo.ID).Updates(stuInfo).Error; err != nil {
		logger.DatabaseLogger.Errorf("更新学生信息失败：%v", err)
		return common.ErrNew(err, common.SysErr)
	}
	return nil
}

// Excelize 获取学生信息并导出为excel
func (*Admin) Excelize() error {
	// 获取所有学生信息
	var data []model.Stu
	if err := model.DB.Model(&model.Stu{}).Preload("Interv").Find(&data).Error; err != nil {
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
		if stu.Message <= 4 {
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
// 这是发送面试结果短息的函数
func sendItvResMsg(pass bool, number string, name string, department string) error {
	client, _err := CreateClient(tea.String(config.Config.AlibabaCloudAccessKeyID), tea.String(config.Config.AlibabaCloudAccessKeySecret))
	if _err != nil {
		return _err
	}

	var tpCode string
	// 短信对应代码,结合自己的短信格式的代码进行设置

	if pass {
		tpCode = "SMS_463215609"
	} else {
		tpCode = "SMS_463195599"
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(number),
		SignName:      tea.String("挑战网"),
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
		_, _err = client.SendSmsWithOptions(sendSmsRequest, runtime)
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

func AliyunSendItvTimeMsg() (fs []FailSend, err error) {
	var intervs []model.Interv
	if err := model.DB.Model(&model.Interv{}).
		Where("time > ? AND time < ?", time.Now().Add(20*time.Minute), time.Now().Add(30*time.Minute)).
		Not("netid = ?", "").Find(&intervs).Error; err != nil {
		return fs, common.ErrNew(errors.New("查询学生信息出错"), common.SysErr)
	}
	for _, interv := range intervs {
		var stu model.Stu
		if interv.NetID == nil || model.DB.Model(&model.Stu{}).Where("netid = ?", interv.NetID).First(&stu).Error != nil {
			fs = append(fs, FailSend{NetID: stu.NetID, ErrCode: -1})
		}
		intervTime := interv.Time
		err := sendItvTimeMsg(stu.Phone, stu.Name, fmt.Sprintf("%d年%d月%d日 %s:%s", intervTime.Year(), intervTime.Month(), intervTime.Day(), intervTime.String()[11:13], intervTime.String()[14:16]), string(stu.Depart))
		if err != nil {
			fs = append(fs, FailSend{NetID: stu.NetID, ErrCode: -1})
		}
	}
	return fs, nil
}

// 发送面试通知短信
func sendItvTimeMsg(phone string, name string, time string, department string) (err error) {
	client, err := CreateClient(&config.Config.AlibabaCloudAccessKeyID, &config.Config.AlibabaCloudAccessKeySecret)
	if err != nil {
		return err
	}
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String("挑战网"),
		TemplateCode:  tea.String("SMS_471010064"),
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

// 发送短信
// 有关数据库的操作
func (a *Admin) AliyunSendItvResMsg() (cocacola interface{}, err error) {
	// 向Aliyun发送请求
	var user []StuMsgInfo
	if err = model.DB.
		Table("students").
		Joins("LEFT JOIN interviews ON students.netid = interviews.netid").
		Select("students.*", "interviews.time").
		Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNew(errors.New("当前学号不存在"), common.NotFoundErr)
		}
		return nil, common.ErrNew(errors.New("没有查到,请联系管理员解决"), common.SysErr)
	}
	for _, value := range user {
		var check bool
		if value.Interv == "已通过" {
			check = true
		} else if value.Interv == "未通过" {
			check = false
		} else {
			continue
		}
		// 发短信在这里
		err := sendItvResMsg(check, value.Phone, value.Name, value.First)
		if err != nil {
			continue
		}
	}
	return nil, nil
}
