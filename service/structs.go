package service

import (
	"time"
	"zhaoxin2025/model"
)

type FailSend struct {
	NetID   string `json:"neteid"`
	ErrCode int    `json:"errcode"`
}

// 已弃用，改为直接使用 Stu
type StuMsgInfo struct {
	Name   string
	Time   time.Time
	Phone  string
	First  string
	Interv string
}

type IntervInfo struct {
	model.Interv
	Name string `json:"name"`
}
