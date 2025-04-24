package model

type AdminLevel string

const (
	Normal AdminLevel = "normal"
	Super  AdminLevel = "super"
)

type Admin struct {
	BaseModel
	NetID    string     `gorm:"column:netid;type:varchar(10);unique;not null;comment:'NetID'"`
	Name     string     `gorm:"column:name;not null;comment:'姓名'"`
	Password string     `gorm:"column:password;not null;comment:'密码'"`
	Level    AdminLevel `gorm:"column:level;not null;type:enum('normal','super');default:'normal';comment:'权限'"`
}
