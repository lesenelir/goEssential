package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"goEssential/common"
	"goEssential/dto"
	"goEssential/model"
	"goEssential/response"
	"goEssential/util"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// http处理函数
func Register(context *gin.Context) {

	DB := common.GetDB()

	// 获取参数
	name := context.PostForm("name")
	telephone := context.PostForm("telephone")
	password := context.PostForm("password")
	// 数据验证
	if len(telephone) != 11 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "手机号必须是11位")
		//context.JSON(http.StatusUnprocessableEntity, gin.H{
		//	"code": 422,
		//	"msg":  "手机号必须是11位",
		//})
		return
	}

	if len(password) < 6 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "密码不能小于6位")
		//context.JSON(http.StatusUnprocessableEntity, gin.H{
		//	"code": 422,
		//	"msg":  "密码不能小于6位",
		//})
		return
	}

	if len(name) == 0 { // 名字没有传值，随机初始化10位字符串
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)

	// 判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户已存在")
		//context.JSON(http.StatusUnprocessableEntity, gin.H{
		//	"code": 422,
		//	"msg":  "用户已存在",
		//})
		return
	}

	// 创建用户
	// 加密用户密码
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil { // 加密错误500
		response.Response(context, http.StatusUnprocessableEntity, 500, nil, "加密错误")
		//context.JSON(http.StatusInternalServerError, gin.H{
		//	"code": 500,
		//	"msg":  "加密错误",
		//})
		return
	}

	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)

	// 发送token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error: %v\n", err)
		return
	}

	// 返回结果
	//context.JSON(http.StatusOK, gin.H{
	//	"code": 200,
	//	"msg":  "注册成功",
	//})
	response.Success(context, gin.H{"token": token}, "注册成功")

}

func Login(context *gin.Context) {
	DB := common.GetDB()

	// 获取参数
	telephone := context.PostForm("telephone")
	password := context.PostForm("password")

	// 数据验证
	if len(telephone) != 11 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "手机号必须是11位")
		//context.JSON(http.StatusUnprocessableEntity, gin.H{
		//	"code": 422,
		//	"msg":  "手机号必须是11位",
		//})
		return
	}

	if len(password) < 6 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		//context.JSON(http.StatusUnprocessableEntity, gin.H{
		//	"code": 422,
		//	"msg":  "密码不能小于6位",
		//})
		return
	}

	// 判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户已存在")
		//context.JSON(http.StatusUnprocessableEntity, gin.H{
		//	"code": 422,
		//	"msg":  "用户已存在",
		//})
		return
	}

	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		//context.JSON(http.StatusBadRequest, gin.H{
		//	"code": 400,
		//	"msg":  "密码错误",
		//})
		response.Response(context, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}

	// 发放token给前端
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "系统异常")
		//context.JSON(http.StatusInternalServerError, gin.H{
		//	"code": 500,
		//	"msg":  "系统异常",
		//})
		log.Printf("token generate err : %v\n", err)
	}

	// 返回结果
	//context.JSON(http.StatusOK, gin.H{
	//	"code": 200,
	//	"data": gin.H{"token": token},
	//	"msg":  "登陆成功",
	//})
	response.Success(context, gin.H{"token": token}, "登陆成功")

}

func Info(context *gin.Context) {
	user, _ := context.Get("user")

	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": dto.ToUserDto(user.(model.User))},
	})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
