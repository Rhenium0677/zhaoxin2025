package model

type Que struct {
	BaseModel
	Question   string     `gorm:"column:question;not null;comment:'问题'"`
	Department Department `gorm:"column:department;enum:('tech','video','art');not null;comment:'所属部门'"`
}
