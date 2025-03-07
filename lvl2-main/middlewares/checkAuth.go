package middlewares

import (
	"fmt"
	"library-management1/database"
	"library-management1/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is misssing"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authToken := strings.Split(authHeader, " ")
		if len(authToken) != 2 || authToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token format"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := authToken[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("auth-api-jwt-secret"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token Expired"})
			c.AbortWithStatus((http.StatusUnauthorized))
			return
		}

		var user models.User
		database.DB.Where("ID = ?", claims["id"]).Find(&user)

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("currentUser", user)
		c.Next()
	}
}
