package services

import (
	"awesomeblog/core"
	"awesomeblog/datamodels"
	"github.com/jinzhu/gorm"
	"time"
)

func GetCommentListForPage(pageIndex, pageNum, article_id int) ([]datamodels.Comment, error) {
	comments := make([]datamodels.Comment, 0)
	db := core.GetInstance()

	var e error
	if e = db.Limit(pageNum).Offset((pageIndex-1)*pageNum).Order("created", true).Where("article_id=?", article_id).Find(&comments).Error; e != nil {
		if e == gorm.ErrRecordNotFound {
			return comments, NoData
		}
	}
	return comments, e
}

func GetCommentTotalByArticle(article_id int) (int, error) {
	db := core.GetInstance()
	count := 0
	if e := db.Model(datamodels.Comment{}).Where("article_id=?", article_id).Count(&count).Error; e != nil {
		return count, e
	}
	return count, nil
}

func AddComment(article_id int, email, content, ip string) error {
	db := core.GetInstance()
	tx := db.Begin()
	comment := datamodels.Comment{ArticleId: article_id, Email: email, Content: content, Ip: ip, Created: time.Now()}
	if e := db.Create(&comment).Error; e != nil {
		tx.Rollback()
		return e
	}
	if e := UpdateCommentNum(article_id); e != nil {
		tx.Rollback()
		return e
	}

	tx.Commit()
	return nil

}
