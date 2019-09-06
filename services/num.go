package services

import (
	"awesomeblog/core"
	"awesomeblog/datamodels"
	"fmt"
	"github.com/jinzhu/gorm"
)

func GetNumByArticleId(article_id int) (datamodels.Num, error) {
	var num datamodels.Num
	db := core.GetInstance()

	var e error
	if e = db.First(&num, article_id).Error; e != nil {
		if e == gorm.ErrRecordNotFound {
			return num, nil
		}
	}
	return num, e
}

func GetNumByArticleIds(article_ids []int) ([]datamodels.Num, error) {
	nums := make([]datamodels.Num, 0)
	db := core.GetInstance()
	if e := db.Where(article_ids).Find(&nums).Error; e != nil {
		return nums, e
	}
	return nums, nil
}

func UpdateCommentNum(article_id int) error {
	db := core.GetInstance()
	var num datamodels.Num
	if e := db.Where(datamodels.Num{ArticleId: article_id}).Attrs(datamodels.Num{CommentNum: 0}).FirstOrCreate(&num).Error; e != nil {
		return e
	}
	fmt.Println(num)
	return db.Model(&num).Update("comment_num", num.CommentNum+1).Error
}

func UpdateReadNum(article_id int) error {
	db := core.GetInstance()
	var num datamodels.Num
	if e := db.Where(datamodels.Num{ArticleId: article_id}).Attrs(datamodels.Num{ReadNum: 0}).FirstOrCreate(&num).Error; e != nil {
		return e
	}

	// num.ReadNum = num.CommentNum + 1
	return db.Model(&num).Update("read_num", num.ReadNum+1).Error
}
