package model

type Que struct {
	BaseModel
	Question   string     `gorm:"column:question;not null;unique;comment:'问题'"`
	Department Department `gorm:"column:department;enum:('tech','video','art');not null;comment:'所属部门'"`
	Url        string     `gorm:"column:url;not null;comment:'问题链接'"`
	Times      int        `gorm:"column:time;not null;default:0;comment:'抽取到的次数'"`
}
