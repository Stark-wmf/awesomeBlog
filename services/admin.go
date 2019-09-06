package services

import (
	"awesomeblog/core"
	"awesomeblog/datamodels"
	"github.com/jinzhu/gorm"
)

func QueryDataByUserName(username string) (datamodels.Admin, error) {
	var admin datamodels.Admin
	db := core.GetInstance()

	var e error
	if e = db.Where("username = ?", username).First(&admin).Error; e != nil {
		if e == gorm.ErrRecordNotFound {
			return admin, NoData
		}
	}
	return admin, e
}

func RegisteUser(username string,password string,nickname string) (error) {

		db := core.GetInstance()
		admin := datamodels.Admin{UserName: username, PassWord:password,GroupId:1,Status:0,NickName:nickname}
		tx := db.Begin()
		if e := tx.Create(&admin).Error; e != nil {
		tx.Rollback()
		return e
	}

		tx.Commit()
		return nil
	}

func EditUser(admin datamodels.Admin) error {
	db := core.GetInstance()
	tx := db.Begin()
	if e := tx.Model(&admin).Updates(datamodels.Admin{
	Status:1}).Error; e != nil {
		tx.Rollback()
		return e
	}

	tx.Commit()
	return nil
}