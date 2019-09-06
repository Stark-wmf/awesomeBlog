package datamodels

type Admin struct {
	Id       int    `gorm:"primary_key" `
	UserName string `gorm:"column:username"`
	NickName string `gorm:"column:nickname"`
	PassWord string `gorm:"column:password"`
	GroupId  int    `gorm:"column:groupId"`
	Status  int    `gorm:"column:status"`
}

func (Admin) TableName() string {
	return "tb_admin"
}
