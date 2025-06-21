package model

type Que struct {
	BaseModel
	Question   string     `gorm:"column:question;not null;unique;comment:'问题'" json:"question"`
	Department Department `gorm:"column:department;enum:('tech','video','art');not null;comment:'所属部门'" json:"department"`
	Url        string     `gorm:"column:url;not null;comment:'问题链接'" json:"url"`
	Times      int        `gorm:"column:times;not null;default:0;comment:'抽取到的次数'" json:"times"`
	// 硬件随机抽取
}
