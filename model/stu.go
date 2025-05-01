package model

import "time"

type Stu struct {
	BaseModel
	OpenID      string     `gorm:"column:openid;type:varchar(50);unique;not null;comment:'微信OpenID'"`
	NetID       string     `gorm:"column:netid;type:varchar(10);unique;not null;comment:'NetID'"`
	Name        string     `gorm:"column:name;not null;comment:'姓名'"`
	Phone       string     `gorm:"column:phone;not null;comment:'电话号码'"`
	School      string     `gorm:"column:school;not null;comment:'书院'"`
	First       Department `gorm:"column:first;enum('tech','video','art');not null;comment:'第一意向'"`
	Second      Department `gorm:"column:second;enum('tech','video','art');not null;comment:'第二意向'"`
	Que         Que        `gorm:"foreignKey:ID;comment:'抽到的问题'"`
	QueTime     time.Time  `gorm:"column:que_time;not null;default:CURRENT_TIMESTAMP;comment:'抽到问题的时间'"`
	Interv      Interv     `gorm:"foreignKey:ID;comment:'预约的面试'"`
	Interviewer string     `gorm:"column:interviewer;not null;comment:'面试官'"`
	Evaluation  string     `gorm:"column:evaluation;not null;comment:'评价'"`
	Star        int        `gorm:"column:star;not null;default:0;comment:'星级'"`
	Message     int        `gorm:"column:messager;not null;default:0;comment:'消息'"`
}

// Message 定义为
// 1.是否在报名后发送订阅消息
// 2.是否在面试时间临近时发送订阅消息
// 3.是否在面试结束后发送结果的订阅消息
// 三个0或1从低位起组成一个二进制数，转为十进制数存储在数据库中
// 例如：000表示不发送任何消息，111表示发送所有消息，001表示只发送报名后的消息
