package model

type Stu struct {
	BaseModel
	OpenID    string     `gorm:"column:openid;type:varchar(50);unique;comment:'微信OpenID'" json:"openid"`
	NetID     string     `gorm:"column:netid;type:varchar(50);unique;comment:'NetID'" json:"netid"`
	Name      string     `gorm:"column:name;not null;comment:'姓名'" json:"name"`
	Phone     string     `gorm:"column:phone;not null;comment:'电话号码'" json:"phone"`
	School    string     `gorm:"column:school;not null;comment:'书院'" json:"school"`
	WhereKnow string     `gorm:"column:whereknow;not null;comment:'在哪里知道的'" json:"whereknow"`
	Mastered  string     `gorm:"column:mastered;not null;comment:'已经会的技能'" json:"mastered"`
	ToMaster  string     `gorm:"column:tomaster;not null;comment:'想要掌握的技能'" json:"tomaster"`
	Depart    Department `gorm:"column:depart;type:enum('tech','video','art','none');default:'none';comment:'主选部门'" json:"depart"`
	Tag       string     `gorm:"column:tag;not null;comment:'标签'" json:"tag"`
	Interv    *Interv    `gorm:"foreignKey:NetID;references:NetID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"interv"`
	Message   int        `gorm:"column:message;not null;default:0;comment:'消息'" json:"message"`
	QueID     int        `gorm:"column:queid;comment:'问题ID'" json:"queid"`
	Work      string     `gorm:"column:work;default:'';comment:'作品链接'" json:"work"`
}

// Message 定义为
// 1.是否在面试结束后发送结果的订阅消息
// 2.是否在面试时间临近时发送订阅消息
// 3.是否在报名面试后发送订阅消息
// 三个0或1从低位起组成一个二进制数，转为十进制数存储在数据库中
// 例如：0b000表示不发送任何消息，0b100表示只发送面试结果的消息，0b001表示只发送报名面试后的消息
