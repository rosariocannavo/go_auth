package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rosariocannavo/go_auth/internal/db"
	"github.com/rosariocannavo/go_auth/internal/models"
	"github.com/rosariocannavo/go_auth/internal/repositories"
	"github.com/rosariocannavo/go_auth/internal/utils"
)

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
		hashedPwd, err := utils.HashPassword(userForm.Password)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		//generate nonce for metamask sign auth
		nonce, err := utils.GenerateRandomNonce()

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
