package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	auth "github.com/rosariocannavo/go_auth/config"
	"github.com/rosariocannavo/go_auth/internal/db"
	"github.com/rosariocannavo/go_auth/internal/models"
	"github.com/rosariocannavo/go_auth/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(c *gin.Context) {

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userRepo := repositories.NewUserRepository(db.Client)

	retrievedUser, err := userRepo.FindUser(&user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not present"})
		return
	}

	fmt.Println()

	errf := bcrypt.CompareHashAndPassword([]byte(retrievedUser.Password), []byte(user.Password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token upon successful authentication
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(auth.SecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.Set("username", user.Username)

	// Return the generated JWT token to the client
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
