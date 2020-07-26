package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"goEssential/model"
)

var DB *gorm.DB

// 数据库初始函数
func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "goEssential"
	username := "root"
	password := "19970122"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}

	// 自动创建数据表
	db.AutoMigrate(&model.User{})
	DB = db
	return db
}


func GetDB() *gorm.DB {
	return DB
}