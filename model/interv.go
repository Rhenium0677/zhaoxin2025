package model

import (
	"time"
)

type Interv struct {
	BaseModel
	NetID       *string    `gorm:"column:netid;type:varchar(50);comment:'NetID'" json:"netid"`
	Time        time.Time  `gorm:"column:time;not null;comment:'面试时间'" json:"time"`
	Interviewer string     `gorm:"column:interviewer;not null;comment:'面试官姓名'" json:"interviewer"`
	Department  Department `gorm:"column:department;not null;comment:'面试部门'" json:"department"`
	Star        int        `gorm:"column:star;not null;default:0;comment:'星级'" json:"star"`
	Evaluation  string     `gorm:"column:evaluation;not null;comment:'评价'" json:"evaluation"`
	Pass        int        `gorm:"column:pass;default:0;not null;comment:'是否通过面试'" json:"pass"`
}
