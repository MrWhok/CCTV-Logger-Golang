package service

import (
	"CCTV-Logger-Golang/src/internal/api/repository"
	"CCTV-Logger-Golang/src/internal/config"
	"CCTV-Logger-Golang/src/internal/pkg/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(c *gin.Context, user models.User)
	Login(c *gin.Context)
	GetAllUsers(c *gin.Context)
	FindUserByUsername(username string) (*models.User, error)
	DeleteUser(username string) error
}

type userService struct {
	repo repository.UserRepository
	cfg  config.Config
}

func NewUserService(repo repository.UserRepository, cfg config.Config) UserService {
	return &userService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *userService) FindUserByUsername(username string) (*models.User, error) {
	return s.repo.FindUserByUsername(username)
}

func (s *userService) Register(c *gin.Context, user models.User) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = string(hashedPassword)
	user.ID = primitive.NewObjectID()

	if err := s.repo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (s *userService) Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.repo.FindUserByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.ID.Hex(),
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // Set expiration time to 1 hour from now
	})

	tokenString, err := token.SignedString([]byte(s.cfg.SecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (s *userService) GetAllUsers(c *gin.Context) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (s *userService) DeleteUser(username string) error {
	return s.repo.DeleteUser(username)
}
