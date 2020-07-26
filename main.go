package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20); not null"`
	Telephone string `gorm:"varchar(110); not null; unique"`
	Password  string `gorm:"size:255; not null"`
}

func main() {
	db := InitDB()
	defer db.Close()

	r := gin.Default()
	r.POST("/api/auth/register", func(context *gin.Context) {
		// 获取参数
		name := context.PostForm("name")
		telephone := context.PostForm("telephone")
		password := context.PostForm("password")
		// 数据验证
		if len(telephone) != 11 {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "手机号必须是11位",
			})
			return
		}

		if len(password) < 6 {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "密码不能小于6位",
			})
			return
		}

		if len(name) == 0 { // 名字没有传值，随机初始化10位字符串
			name = RandomString(10)
		}

		log.Println(name, telephone, password)

		// 判断手机号是否存在
		if isTelephoneExist(db, telephone) {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "用户名存在",
			})
		}

		// 创建用户
		newUser := User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)

		// 返回结果

		context.JSON(http.StatusOK, gin.H{
			"msg": "注册成功",
		})
	})

	panic(r.Run())
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func RandomString(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

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
	db.AutoMigrate(&User{})
	return db
}
