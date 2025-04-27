package model

import (
	"time"
)

type Interv struct {
	BaseModel
	Time time.Time `gorm:"column:time;not null;comment:'面试时间'"`
}
