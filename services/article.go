package services

import (
	"awesomeblog/core"
	"awesomeblog/datamodels"
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
)

type AST int8

const (
	_ AST = iota - 2
	ArticleAll
	ArticleNORMAL
	ArticleDELETE
)

func GetArticleList(num int, status AST) ([]datamodels.Article, error) {
	var articles []datamodels.Article
	var e error
	db := core.GetInstance()

	if num <= 0 {
		if status == ArticleAll {
			e = db.Model(&datamodels.Article{}).Order("updated_at desc").Find(&articles).Error
		} else {
			e = db.Model(&datamodels.Article{}).Order("updated_at desc").Where("status=?", status).Find(&articles).Error
		}
	} else {
		if status == ArticleAll {
			e = db.Model(&datamodels.Article{}).Order("updated_at desc").Limit(num).Find(&articles).Error
		} else {
			e = db.Model(&datamodels.Article{}).Order("updated_at desc").Limit(num).Where("status=?", status).Find(&articles).Error
		}
	}
	if e != nil {
		return articles, e
	}

	categorys, e := GetCategoryList(CategoryNORMAL)
	if e != nil {
		return articles, e
	}

	tmpCategorys := make(map[int]string)
	for i := 0; i < len(categorys); i++ {
		tmpCategorys[categorys[i].ID] = categorys[i].Name
	}

	for i := 0; i < len(articles); i++ {
		articles[i].CategoryName = tmpCategorys[articles[i].CategoryId]
	}
	return articles, nil
}

func AddArticle(category_id int, title string, content string, top int8, depiction string, image string) error {

	db := core.GetInstance()
	article := datamodels.Article{CategoryId: category_id, Title: title, Content: content, Top: top, Depiction: depiction, Image: image}
	tx := db.Begin()
	if top == 2 {
		if e := tx.Model(&datamodels.Article{}).Updates(map[string]int8{"top": 1}).Error; e != nil {
			tx.Rollback()
			return e
		}
	}

	if e := tx.Create(&article).Error; e != nil {
		tx.Rollback()
		return e
	}

	tx.Commit()
	return nil
}

func EditArticle(article datamodels.Article) error {
	db := core.GetInstance()
	tx := db.Begin()
	if article.Top == 2 {
		if e := tx.Model(&datamodels.Article{}).Updates(map[string]int8{"top": 1}).Error; e != nil {
			tx.Rollback()
			return e
		}
	}
	if e := tx.Model(&article).Updates(datamodels.Article{
		CategoryId: article.CategoryId,
		Title:      article.Title,
		Content:    article.Content,
		Top:        article.Top,
		Depiction:  article.Depiction,
		Image:      article.Image}).Error; e != nil {
		tx.Rollback()
		return e
	}

	tx.Commit()
	return nil
}

func GetArticleCountByCategoryId(category_id int) (int, error) {
	db := core.GetInstance()

	var count int
	e := db.Model(datamodels.Article{}).Where("category_id = ? and status = 0", category_id).Count(&count).Error
	if e != nil {
		return 0, e
	}
	return count, nil
}

func DeleteArticle(id int) (bool, error) {

	db := core.GetInstance()
	e := db.Model(&datamodels.Article{}).Where("id = ?", id).Update(map[string]int{"status": 1}).Error
	if e != nil {
		return false, e
	}
	return true, nil
}

func GetArticleById(id int) (datamodels.Article, error) {
	var article datamodels.Article
	db := core.GetInstance()

	if e := db.First(&article, id).Error; e != nil {
		if e == gorm.ErrRecordNotFound {
			return article, nil
		}
		return article, e
	}
	return article, nil
}

func GetTopArticle() (datamodels.Article, error) {
	var article datamodels.Article
	db := core.GetInstance()

	if e := db.Model(&datamodels.Article{}).Where("status=? ", 0).Where("top=?", 2).First(&article).Error; e != nil {
		if e == gorm.ErrRecordNotFound {
			return article, nil
		}
		return article, e
	}
	return article, nil

}

func GetArticleListForPage(pageIndex, pageNum, category_id int) ([]datamodels.Article, error) {
	articles := make([]datamodels.Article, pageNum)
	db := core.GetInstance()
	var e error
	if category_id == 0 {
		e = db.Limit(pageNum).Offset((pageIndex-1)*pageNum).Order("created_at", true).Where("status=?", ArticleNORMAL).Find(&articles).Error
	} else {
		e = db.Limit(pageNum).Offset((pageIndex-1)*pageNum).Order("created_at", true).Where("status=?", ArticleNORMAL).Where("category_id=?", category_id).Find(&articles).Error
	}

	if e != nil {
		if e == sql.ErrNoRows {
			return articles, nil
		} else {
			return articles, e
		}
	}
	categorys, e := GetCategoryList(CategoryNORMAL)
	if e != nil {
		return articles, e
	}

	tmpCategorys := make(map[int]string)
	for i := 0; i < len(categorys); i++ {
		tmpCategorys[categorys[i].ID] = categorys[i].Name
	}

	for i := 0; i < len(articles); i++ {
		articles[i].CategoryName = tmpCategorys[articles[i].CategoryId]
	}

	return articles, nil
}

func GetArticleTotal() (int, error) {
	db := core.GetInstance()
	count := 1
	if e := db.Model(datamodels.Article{}).Where("status=?", 0).Count(&count).Error; e != nil {
		return count, e
	}
	return count, nil
}

func GetArticleByTwo(article_id int) ([]datamodels.Article, error) {

	results := make([]datamodels.Article, 0)
	db := core.GetInstance()
	rows, _ := db.Raw(fmt.Sprintf("select id, title from tb_article where id in (select case when SIGN(id-14)>0 THEN MIN(id) when SIGN(id-%d)<0 THEN MAX(id) ELSE id end from tb_article where status = 0 GROUP BY SIGN(id-%d) ORDER BY SIGN(id-%d) ) ORDER BY id ", article_id, article_id, article_id)).Rows() // (*sql.Rows, error)
	defer rows.Close()
	var (
		id    int
		title string
	)
	for rows.Next() {
		rows.Scan(&id, &title)
		var t datamodels.Article
		if id == article_id {
			continue
		}
		if id < article_id {
			t = datamodels.Article{BaseModel: datamodels.BaseModel{ID: id}, Title: title, Content: "s"}
		} else {
			t = datamodels.Article{BaseModel: datamodels.BaseModel{ID: id}, Title: title, Content: "x"}
		}
		results = append(results, t)
	}
	return results, nil
}
