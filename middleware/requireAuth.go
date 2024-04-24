package middleware

import (
	"JWT-AUTH-GIN/initializers"
	"JWT-AUTH-GIN/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"os"
	"time"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString, err := c.Cookie("Authorization")

		// Check if the token is missing
		if err != nil || tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			log.Fatal(err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			var user models.User
			initializers.DB.First(&user, claims["sub"])
			if user.ID == 0 {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			c.Set("user", user)
			c.Next()

		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}
}
