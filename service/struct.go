package service

import (
	"time"
	"zhaoxin2025/common"
	"zhaoxin2025/model"
)

// Admin 相关结构体
type SelectStu struct {
	NetID       string           `json:"netid" binding:"omitempty,len=10,numeric"`
	Name        string           `json:"name" binding:"omitempty"`
	Phone       string           `json:"phone" binding:"omitempty,len=11,numeric"`
	School      string           `json:"school" binding:"omitempty"`
	First       model.Department `json:"first" binding:"omitempty,oneof=tech video art"`
	Second      model.Department `json:"second" binding:"omitempty,oneof=tech video art"`
	Interviewer string           `json:"interviewer" binding:"omitempty"`
	Star        int              `json:"star" binding:"omitempty"`
}
type UpdateStu struct {
	NetID       string           `json:"netid" binding:"required,len=10,numeric"`
	Name        string           `json:"name" binding:"omitempty"`
	Phone       string           `json:"phone" binding:"omitempty"`
	School      string           `json:"school" binding:"omitempty"`
	Mastered    string           `json:"mastered" binding:"omitempty"`
	ToMaster    string           `json:"tomaster" binding:"omitempty"`
	First       model.Department `json:"first" binding:"omitempty,oneof=tech video art"`
	Second      model.Department `json:"second" binding:"omitempty,oneof=tech video art"`
	QueID       int              `json:"que_id" binding:"omitempty,numeric"`
	QueTime     time.Time        `json:"que_time" binding:"omitempty"`
	Interv      bool             `json:"interv" binding:"omitempty"`
	Interviewer string           `json:"interviewer" binding:"omitempty"`
	Evaluation  string           `json:"evaluation" binding:"omitempty"`
	Star        int              `json:"star" binding:"omitempty"`
}

// Stu 相关结构体
type StuUpdate struct {
	NetID    string `json:"netid" binding:"required,len=10,numeric"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required,len=11,numeric"`
	Mail     string `json:"mail" binding:"required,email"`
	School   string `json:"school" binding:"required"`
	Mastered string `json:"mastered" binding:"required"`
	ToMaster string `json:"tomaster" binding:"required"`
}

// Interv 相关结构体
type IntervUpdate struct {
	ID          int              `json:"id" binding:"required"`
	Interviewer string           `json:"interviewer" binding:"omitempty"`
	Time        time.Time        `json:"time" binding:"omitempty"`
	NetID       string           `json:"netid" binding:"omitempty,len=10,numeric"`
	Department  model.Department `json:"department" binding:"omitempty,oneof=tech video art"`
	Star        int              `json:"star" binding:"omitempty"`
	Evaluation  string           `json:"evaluation" binding:"omitempty"`
	Pass        string           `json:"pass" binding:"omitempty,oneof=true false"`
}
type GetInterv struct {
	ID          int              `json:"id" binding:"omitempty"`
	Department  model.Department `json:"department,omitempty"`
	Interviewer string           `json:"interviewer,omitempty"`
	Pass        string           `json:"pass" binding:"omitempty,oneof=true false"`
	Date        time.Time        `json:"date,omitempty"`
	common.PagerForm
}

// Que 相关结构体
type UpdateQue struct {
	ID         int              `json:"id" binding:"required"`
	QueID      int              `json:"queid" binding:"omitempty"`
	Department model.Department `json:"department" binding:"omitempty,oneof=tech video art"`
	Times      int              `json:"times" binding:"omitempty"`
}
