package router

import (
	"CCTV-Logger-Golang/src/internal/api/handler"
	"CCTV-Logger-Golang/src/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterImageRoutes(r *gin.Engine, imageHandler handler.ImageHandler) {
	imageGroup := r.Group("/images")
	{
		imageGroup.GET("/", middleware.Authenticate, imageHandler.GetAllImages)
		imageGroup.GET("/:id", middleware.Authenticate, imageHandler.FindImage, imageHandler.GetImage)
		imageGroup.POST("/upload", middleware.Authenticate, imageHandler.UploadImage)
		imageGroup.DELETE("/delete/:id", middleware.Authenticate, imageHandler.FindImage, imageHandler.DeleteImage)
		imageGroup.PATCH("/favorite/:id", middleware.Authenticate, imageHandler.FindImage, imageHandler.FavoriteImage)
	}
}
