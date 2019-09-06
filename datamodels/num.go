package datamodels

type Num struct {
	ArticleId  int `gorm:"primary_key"`
	CommentNum int `gorm:"column:comment_num"`
	ReadNum    int `gorm:"column:read_num"`
}

func (Num) TableName() string {
	return "tb_num"
}
