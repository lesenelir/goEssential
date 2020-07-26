package main

import (
	"github.com/gin-gonic/gin"
	"goEssential/controller"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	return r
}