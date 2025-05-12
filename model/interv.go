package model

import (
	"time"
)

type Interv struct {
	BaseModel
	NetID       *string    `gorm:"column:netid;type:varchar(10);comment:'NetID'"`
	Time        time.Time  `gorm:"column:time;not null;comment:'面试时间'"`
	Interviewer string     `gorm:"column:interviewer;not null;comment:'面试官姓名'"`
	Department  Department `gorm:"column:department;not null;comment:'面试部门'"`
	QueID       int        `gorm:"column:queid;comment:'问题ID'"`
	Que         *Que       `gorm:"foreignKey:QueID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;comment:'问题'"`
	Star        int        `gorm:"column:star;not null;default:0;comment:'星级'"`
	Evaluation  string     `gorm:"column:evaluation;not null;comment:'评价'"`
	Pass        bool       `gorm:"column:pass;not null;comment:'是否通过面试'"`
}
