package middleware

import (
	"github.com/gin-gonic/gin"
	"goEssential/common"
	"goEssential/model"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 获取authorization header
		tokenString := context.GetHeader("Authorization")

		// validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			context.Abort()
			return
		}

		// 提取token有效部分
		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)

		if err != nil || !token.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			return
		}

		// 验证通过后获取claim中的userId
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		// 用户
		if user.ID == 0 {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			context.Abort()
			return
		}

		// 用户存在 将user信息写入上下文
		context.Set("user", user)

		context.Next()

	}
}
