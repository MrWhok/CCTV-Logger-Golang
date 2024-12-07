package handler

import (
	"CCTV-Logger-Golang/src/internal/api/service"
	"CCTV-Logger-Golang/src/internal/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	ValidateUserData(ctx *gin.Context)
	CheckUserExists(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) Login(ctx *gin.Context) {
	h.userService.Login(ctx)
}

func (h *userHandler) GetAllUsers(ctx *gin.Context) {
	h.userService.GetAllUsers(ctx)
}

func (h *userHandler) DeleteUser(ctx *gin.Context) {
	username := ctx.Param("username")

	err := h.userService.DeleteUser(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting user", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *userHandler) ValidateUserData(ctx *gin.Context) {
	var user models.UserData

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data", "error": err.Error()})
		ctx.Abort()
		return
	}

	// Store the user data in the context
	ctx.Set("user", user)

	ctx.Next()
}

func (h *userHandler) CheckUserExists(ctx *gin.Context) {
	// Retrieve the user data from the context
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "User data not found"})
		ctx.Abort()
		return
	}

	userData := user.(models.UserData)

	userExists, err := h.userService.FindUserByUsername(userData.Username)
	if err == nil && userExists != nil {
		ctx.JSON(http.StatusConflict, gin.H{"message": "User already exists"})
		ctx.Abort()
		return
	}

	ctx.Next()
}

func (h *userHandler) Register(ctx *gin.Context) {
	// Retrieve the user data from the context
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "User data not found"})
		return
	}

	userData := user.(models.UserData)

	// Convert the user data to the models.User type
	newUser := models.User{
		Username:  userData.Username,
		Password:  userData.Password,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
	}

	h.userService.Register(ctx, newUser)
}
