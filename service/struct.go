package service

import "zhaoxin2025/model"

type SelectStu struct {
	ID          int              `json:"id" binding:"omitempty"`
	NetID       string           `json:"netid" binding:"omitempty,len=10,numeric"`
	Name        string           `json:"name" binding:"omitempty"`
	Phone       string           `json:"phone" binding:"omitempty"`
	School      string           `json:"school" binding:"omitempty"`
	First       model.Department `json:"first" binding:"omitempty,oneof=tech video art"`
	Second      model.Department `json:"second" binding:"omitempty,oneof=tech video art"`
	Interviewer string           `json:"interviewer" binding:"omitempty"`
	Star        int              `json:"star" binding:"omitempty"`
}
