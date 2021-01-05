package model

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

func ConnMysql(host, port, user, pass, dbName string) error {
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, host, port, dbName))
	if err != nil {
		fmt.Println(err)
		log.Fatal("mysql conn error")
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(8)
	db.DB().SetMaxOpenConns(10)
	db.LogMode(true)
	return err
}
