package service

import (
	"zhaoxin2025/common"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"
)

type Stu struct{}

func (*Stu) Login(netid string, code string) (string, string, error) {
	// openid, sessionKey := WxLogin(code)
	// var info model.Stu
	// if err := model.DB.Where("netid = ?", netid).FirstOrCreate(&info).Error; err != nil {
	// 	logger.DatabaseLogger.Errorf("查询学生信息失败: %v", err)
	// 	return "", "", common.ErrNew(err, common.SysErr)
	// }
	// if info.OpenID != openid {
	// 	return "", "", common.ErrNew(errors.New("数据库中openid与用code获取到的openid不一致"), common.AuthErr)
	// }
	// token :=
	// return token, openid, nil
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
