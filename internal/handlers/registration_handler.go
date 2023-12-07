package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

func generateRandomNonce() (string, error) {
	// Define the length of the nonce
	nonceLength := 16 // You can adjust the length as needed

	// Create a byte slice to store the random nonce
	nonce := make([]byte, nonceLength)

	// Read random bytes into the nonce slice
	_, err := rand.Read(nonce)
	if err != nil {
		return "", err
	}

	fmt.Println(nonce)
	// Encode the random bytes to a hexadecimal string
	nonceString := "0x" + hex.EncodeToString(nonce)

	fmt.Println("generated nonce", nonceString+"\n")

	return nonceString, nil
}

func HandleRegistration(c *gin.Context) {
	userRepo := repositories.NewUserRepository(db.Client)

	var userForm models.UserForm
	//retrieve the partial user information from form
	if err := c.BindJSON(&userForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	//check if the user is already registered
	isPresent, err := userRepo.CheckIfUserIsPresent(userForm.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching database"})
		return
	}

	if isPresent {
		c.JSON(http.StatusForbidden, gin.H{"error": "User already present"})
		return
	} else {

		//if user is not present
		//hash his psw and store him in the db
		//give him a role
		//give him a nonce

		var user models.User

		//hash the user password
		hashedPwd, err := hashPassword(userForm.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		//generate nonce for metamask sign auth
		nonce, err := generateRandomNonce()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": " bad nonce generation"})
			return
		}

		user.Username = userForm.Username
		user.Password = hashedPwd
		user.MetamaskAddress = userForm.MetamaskAddress
		user.Nonce = nonce

		userRepo.CreateUser(&user)

		c.JSON(http.StatusOK, gin.H{"message": "user created succesfully"})
	}
}
