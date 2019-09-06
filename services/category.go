package services

import (
	"awesomeblog/core"
	"awesomeblog/datamodels"
	"fmt"
	"github.com/jinzhu/gorm"
)

type CST int8

const (
	_ CST = iota - 1
	CategoryAll
	CategoryNORMAL
	CategoryDELETE
)

func AddCategory(name string) error {

	db := core.GetInstance()
	category := datamodels.Category{Name: name, Status: 1}
	return db.Create(&category).Error
}

/**
	根据id获取分类信息
params:
	id:分类自增长编号
result:
	Category
	error
*/
func GetCategoryById(id int) (datamodels.Category, error) {
	var category datamodels.Category
	db := core.GetInstance()

	if e := db.First(&category, id).Error; e != nil {
		if e == gorm.ErrRecordNotFound {
			return category, nil
		}
		return category, e
	}
	return category, nil
}

func DeleteCategory(id int) (bool, error) {

	db := core.GetInstance()
	e := db.Model(&datamodels.Category{}).Where("id = ?", id).Update(map[string]int{"status": 2}).Error
	if e != nil {
		return false, e
	}
	return true, nil
}

func EditCategory(id int, name string) (bool, error) {

	db := core.GetInstance()
	var categoryModel datamodels.Category
	if e := db.First(&categoryModel, id).Error; e != nil {
		if e == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, e
	}

	categoryModel.Name = name
	e := db.Model(&categoryModel).Update(&categoryModel).Error
	if e != nil {
		fmt.Println(e.Error())
		return false, e
	}
	return true, nil
}

/**
获取分类列表
status:0 全部 1:正常 2:删除
*/

func GetCategoryList(status CST) ([]datamodels.Category, error) {
	var categorys []datamodels.Category
	db := core.GetInstance()
	var e error
	if status != CategoryAll {
		e = db.Model(&datamodels.Category{}).Where("status = ?", status).Find(&categorys).Error
	} else {
		e = db.Model(&datamodels.Category{}).Find(&categorys).Error
	}
	if e != nil {
		return categorys, e
	}
	return categorys, nil

}
