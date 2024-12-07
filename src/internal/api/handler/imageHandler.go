package handler

import (
	"CCTV-Logger-Golang/src/internal/api/service"
	"CCTV-Logger-Golang/src/internal/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageHandler interface {
	GetAllImages(ctx *gin.Context)
	GetImage(ctx *gin.Context)
	UploadImage(ctx *gin.Context)
	DeleteImage(ctx *gin.Context)
	FavoriteImage(ctx *gin.Context)
	FindImage(ctx *gin.Context)
}

type imageHandler struct {
	imageService service.ImageService
}

func NewImageHandler(imageService service.ImageService) ImageHandler {
	return &imageHandler{
		imageService: imageService,
	}
}

func (h *imageHandler) FindImage(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		ctx.Abort()
		return
	}

	image, err := h.imageService.FindByID(objectId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Image not found"})
		ctx.Abort()
		return
	}

	ctx.Set("image", image)
	ctx.Next()
}

func (h *imageHandler) GetAllImages(ctx *gin.Context) {
	images, err := h.imageService.GetAllImages()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, images)
}

func (h *imageHandler) GetImage(ctx *gin.Context) {
	image := ctx.MustGet("image").(*models.Image)
	response := h.imageService.GetImage(image)
	ctx.JSON(http.StatusOK, response)
}

func (h *imageHandler) UploadImage(ctx *gin.Context) {
	var request struct {
		FilePath    string `json:"filePath"`
		TotalEntity int    `json:"totalEntity"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	if request.FilePath == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "No file path provided"})
		return
	}

	response, err := h.imageService.UploadImage(request.FilePath, request.TotalEntity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (h *imageHandler) DeleteImage(ctx *gin.Context) {
	image := ctx.MustGet("image").(*models.Image)
	response, err := h.imageService.DeleteImage(image)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *imageHandler) FavoriteImage(ctx *gin.Context) {
	var requestBody struct {
		Starred bool `json:"starred"`
	}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid starred value. It must be a boolean."})
		return
	}

	image := ctx.MustGet("image").(*models.Image)
	response, err := h.imageService.FavoriteImage(image, requestBody.Starred)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
