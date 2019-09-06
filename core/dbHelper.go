package core

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var dbInstance *gorm.DB

func InitDB() error {
//	c := GetConfig()
//	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", c.DBUser, c.DBPasswd, c.DBHost, c.DBPort, c.DBName)
    dns:="root:191513@tcp(localhost:3306)/awesomeblog?charset=utf8&parseTime=True&loc=Local"
	var err error
	dbInstance, err = gorm.Open("mysql", dns)

	if err != nil || dbInstance.DB().Ping() != nil {
		return err
	}

	return nil
}

func GetInstance() (conn *gorm.DB) {
	return dbInstance
}

func Close() {
	dbInstance.Close()
}
