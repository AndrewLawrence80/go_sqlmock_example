package main

import (
	"example/entity"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var once sync.Once

func connect() {
	once.Do(func() {
		var err error
		dsn := "root:root@tcp(127.0.0.1:43306)/test?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	})
}

func initDB() {
	db.AutoMigrate(&entity.User{})
}

func main() {
	connect()
	initDB()
}
