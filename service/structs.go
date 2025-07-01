package service

import "time"

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
