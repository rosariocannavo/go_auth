package handlers

import (
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

func hashPassword(password string) (string, error) {
	// Hashing the password with a cost of 14 (adjust as needed)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func HandleRegistration(c *gin.Context) {

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userRepo := repositories.NewUserRepository(db.Client)

	isPresent, err := userRepo.CheckIfUserIsPresent(user.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching database"})
		return
	}

	if isPresent {
		c.JSON(http.StatusForbidden, gin.H{"error": "User already present"})
		return
	} else {
		//if user is not present hash his psw and store him in the db
		hashedPwd, err := hashPassword(user.Password)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		user.Password = hashedPwd
		userRepo.CreateUser(&user)
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
