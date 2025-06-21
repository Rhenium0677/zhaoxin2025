package model

type AdminLevel string

const (
	Normal AdminLevel = "normal"
	Super  AdminLevel = "super"
)

type Department string

const (
	Tech  Department = "tech"
	Video Department = "video"
	Art   Department = "art"
	None  Department = "none"
)

type Admin struct {
	BaseModel
	NetID    string     `gorm:"column:netid;type:varchar(10);unique;not null;comment:'NetID'" json:"netid"`
	Name     string     `gorm:"column:name;not null;comment:'姓名'" json:"name"`
	Password string     `gorm:"column:password;not null;comment:'密码'" json:"password"`
	Level    AdminLevel `gorm:"column:level;not null;type:enum('normal','super');default:'normal';comment:'权限'" json:"level"`
}
