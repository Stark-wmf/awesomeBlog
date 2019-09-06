package datamodels

type Article struct {
	BaseModel
	CategoryId   int    `gorm:"category_id" json:"category_id"`
	CategoryName string `gorm:"-" json:"category_name"`
	Title        string `gorm:"title" json:"title"`
	Content      string `gorm:"content" json:"content"`
	Status       int8   `gorm:"status" json:"status"`
	Top          int8   `gorm:"top" json:"top"`
	Depiction    string `gorm:"depiction" json:"depiction"`
	Image        string `gorm:"image" json:"image"`
	ReadNum      int    `gorm:"-" json:"read_num"`
	CommentNum   int    `gorm:"-" json:"comment_num"`
}

func (Article) TableName() string {
	return "tb_article"
}
