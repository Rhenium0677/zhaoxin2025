package model

import (
	"time"
)

type Interv struct {
	BaseModel
	NetID       *string    `gorm:"column:netid;type:varchar(50);comment:'NetID'" json:"netid"`
	Time        time.Time  `gorm:"column:time;not null;comment:'面试时间'" json:"time"`
	Interviewer string     `gorm:"column:interviewer;not null;comment:'面试官姓名'" json:"interviewer"`
	Department  Department `gorm:"column:department;not null;default:'none';comment:'面试部门'" json:"department"`
	Star        int        `gorm:"column:star;not null;default:0;comment:'星级'" json:"star"`
	Evaluation  string     `gorm:"column:evaluation;not null;comment:'评价'" json:"evaluation"`
	Pass        int        `gorm:"column:pass;default:0;not null;comment:'是否通过面试'" json:"pass"`
	QueID       int        `gorm:"column:queid;comment:'问题ID'" json:"queid"`
	QueTime     int        `gorm:"column:quetime;default:0;not null;comment:'抽题时的时间戳'" json:"quetime"`
	Status      int        `gorm:"column:status;default:0;not null;comment:'状态，0-未开始，1-进行中，2-已结束'" json:"status"`
}
