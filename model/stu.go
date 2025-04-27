package model

type Stu struct {
	BaseModel
	NetID       string     `gorm:"column:netid;type:varchar(10);unique;not null;comment:'NetID'"`
	Name        string     `gorm:"column:name;not null;comment:'姓名'"`
	Phone       string     `gorm:"column:phone;not null;comment:'电话号码'"`
	School      string     `gorm:"column:school;not null;comment:'书院'"`
	First       Department `gorm:"column:first;enum('tech','video','art');not null;comment:'第一意向'"`
	Second      Department `gorm:"column:second;enum('tech','video','art');not null;comment:'第二意向'"`
	Que         Que        `gorm:"foreignKey:ID;comment:'抽到的问题'"`
	Interv      Interv     `gorm:"foreignKey:ID;comment:'预约的面试'"`
	Interviewer string     `gorm:"column:interviewer;not null;comment:'面试官'"`
	Evaluation  string     `gorm:"column:evaluation;not null;comment:'评价'"`
}
