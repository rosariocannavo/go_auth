package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	auth "github.com/rosariocannavo/go_auth/config"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		//fmt.Printf(tokenString + "\n")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return auth.JWTSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("Invalid token claims")
			c.Abort()
			return
		}

		if !ok {
			fmt.Println("Username is not a string or doesn't exist")
			c.Abort()
			return
		}

		//retrieve information based on the token - could do this in other client
		//fmt.Println("username authorized", claims["username"])
		c.Set("claims", claims)

		c.Next()
	}
}
