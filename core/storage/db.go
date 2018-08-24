package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

var dbClient *gorm.DB

func NewDB(userName, passWord, dbHost, dbPort, database string) error {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		userName, passWord, dbHost, dbPort, database)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	db.DB().SetMaxOpenConns(400)
	db.DB().SetMaxIdleConns(300)
	//dbClient = db.Debug()
	dbClient = db
	return nil
}

func StopDB() {
	dbClient.Close()
}
