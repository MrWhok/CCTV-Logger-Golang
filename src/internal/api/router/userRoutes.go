package router

import (
	"CCTV-Logger-Golang/src/internal/api/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, userHandler handler.UserHandler) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", userHandler.ValidateUserData, userHandler.CheckUserExists, userHandler.Register)
		userGroup.POST("/login", userHandler.Login)
		userGroup.GET("/all", userHandler.GetAllUsers)
		userGroup.DELETE("/delete/:username", userHandler.DeleteUser)
	}
}
