package service

import (
	"time"
	"zhaoxin2025/model"
)

type FailSend struct {
	NetID   string `json:"neteid"`
	ErrCode int    `json:"errcode"`
}

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
