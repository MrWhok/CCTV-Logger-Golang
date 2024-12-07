package app

import (
	"CCTV-Logger-Golang/src/internal/api/handler"
	"CCTV-Logger-Golang/src/internal/api/repository"
	"CCTV-Logger-Golang/src/internal/api/router"
	"CCTV-Logger-Golang/src/internal/api/service"
	"CCTV-Logger-Golang/src/internal/config"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	cfg := config.LoadConfig()

	// Body parser
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Serve static files
	r.Static("/staticimages", "./public/staticimages")

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	imageRepo := repository.NewImageRepository()

	// Initialize services
	userService := service.NewUserService(userRepo, cfg)
	imageService := service.NewImageService(imageRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	imageHandler := handler.NewImageHandler(imageService)

	// Register routes
	router.RegisterUserRoutes(r, userHandler)
	router.RegisterImageRoutes(r, imageHandler)

	return r
}
