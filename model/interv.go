package model

import (
	"time"
)

type Interv struct {
	BaseModel
	Time        time.Time  `gorm:"column:time;not null;comment:'面试时间'"`
	Interviewer string     `gorm:"column:interviewer;not null;comment:'面试官姓名'"`
	Department  Department `gorm:"column:department;not null;comment:'面试部门'"`
	Que         Que        `gorm:"foreignKey:ID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Star        int        `gorm:"column:star;not null;default:0;comment:'星级'"`
	Evaluation  string     `gorm:"column:evaluation;not null;comment:'评价'"`
}
