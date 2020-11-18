package db

import (
    "fmt"
    "sync"
    "errors"

    orm "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "github.com/spf13/viper"
)

var once sync.Once
var err error 
var db *orm.DB

// 单例模式
func GetDB() *orm.DB {
    once.Do(func() {
        db = openPool()
    })

    return db
}

func openPool() *orm.DB {
	// 这里有一种常见的拼接字符串的方式
    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", viper.GetString("db.username"), viper.GetString("db.password"), viper.GetString("db.host"), viper.GetString("db.name"), viper.GetString("db.charset"))
    db, err := orm.Open("mysql", dsn)
    if err != nil {
        panic(errors.New("mysql连接失败"))
    }

    // 连接数配置也可以写入配置，在此读取
    db.DB().SetMaxIdleConns(viper.GetInt("db.MaxIdleConns"))
    db.DB().SetMaxOpenConns(viper.GetInt("db.MaxOpenConns"))
    db.LogMode(true)
    return db
}

func CloseDB(db *orm.DB)  {
	db.Close()
}