package datamodels

import "time"

type Comment struct {
	Id        int       `gorm:"primary_key" `
	ArticleId int       `gorm:"column:article_id"`
	Email     string    `gorm:"column:email"`
	Content   string    `gorm:"column:content"`
	Created   time.Time `gorm:"column:created"`
	Ip        string    `gorm:"column:ip"`
}

func (Comment) TableName() string {
	return "tb_comment"
}
