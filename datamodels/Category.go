package datamodels

type Category struct {
	BaseModel
	Name   string `gorm:"name" json:"name"`
	Status int8   `gorm:"status" json:"status"`
}

func (Category) TableName() string {
	return "tb_category"
}
