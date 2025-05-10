package service

import (
	"time"
)

// TimeRange 定义一个时间范围结构体
// 包含开始时间和结束时间
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// DayRange 接收一个时间参数，返回同一天的开始和结束 TimeRange 结构体
func DayRange(t time.Time) TimeRange {
	return TimeRange{
		Start: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()),
		End:   time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location()),
	}
}
