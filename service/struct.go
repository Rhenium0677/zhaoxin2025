package service

import (
	"time"
	"zhaoxin2025/model"
)

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

type UpdateQue struct {
	ID         int              `json:"id" binding:"required"`
	Question   string           `json:"question" binding:"omitempty"`
	Department model.Department `json:"department" binding:"omitempty,oneof=tech video art"`
	Url        string           `json:"url" binding:"omitempty"`
	Times      int              `json:"times" binding:"omitempty"`
}
