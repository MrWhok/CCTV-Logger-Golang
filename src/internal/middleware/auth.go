package middleware

import (
	"CCTV-Logger-Golang/src/internal/config"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	cfg := config.LoadConfig()
	log.Println("Secret Key (Authenticate):", cfg.SecretKey) // Log the secret key

	tokenString := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)
	if tokenString == "" {
		log.Println("Authorization header missing or invalid format")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Access denied"})
		c.Abort()
		return
	}

	log.Println("Token String:", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.SecretKey), nil
	})

	if err != nil {
		log.Println("Error parsing token:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token"})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Println("Claims:", claims)
		c.Set("userId", claims["userId"].(string))
		c.Next()
	} else {
		log.Println("Token is invalid or claims are not valid")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token"})
		c.Abort()
	}
}
