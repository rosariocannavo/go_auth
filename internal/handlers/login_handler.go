package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rosariocannavo/go_auth/config"
	"github.com/rosariocannavo/go_auth/internal/db"
	"github.com/rosariocannavo/go_auth/internal/models"
	"github.com/rosariocannavo/go_auth/internal/repositories"

	"github.com/rosariocannavo/go_auth/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var retrievedUser *models.User

func HandleLogin(c *gin.Context) {
	userRepo := repositories.NewUserRepository(db.Client)
	var userForm models.UserForm

	if err := c.BindJSON(&userForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var errdb error
	retrievedUser, errdb = userRepo.FindUser(userForm.Username)

	fmt.Println("username", userForm.Username)
	if errdb != nil {
		fmt.Println("USER NOT PRESENT")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not present"})
		return
	}

	errf := utils.CompareHashPwd(retrievedUser.Password, userForm.Password)

	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	//add metamask authentication
	c.JSON(http.StatusAccepted, gin.H{"Nonce": retrievedUser.Nonce})

}

func HandleverifySignature(c *gin.Context) {

	userRepo := repositories.NewUserRepository(db.Client)

	var requestData struct {
		Nonce     string `json:"message"`
		Address   string `json:"address"`
		Signature string `json:"signature"`
	}

	// Bind JSON body to struct
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// fmt.Println("mesg " + requestData.Nonce)
	// fmt.Println("addr " + requestData.Address)
	// fmt.Println("sig " + requestData.Signature + "\n")

	isSignatureVerified := utils.CheckSig(requestData.Address, requestData.Signature, []byte(requestData.Nonce))

	if isSignatureVerified {
		// Generate JWT token upon successful authentication
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":              retrievedUser.ID,
			"username":        retrievedUser.Username,
			"metamaskAddress": retrievedUser.MetamaskAddress,
			"nonce":           retrievedUser.Nonce,
			"role":            retrievedUser.Role,

			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString(config.JWTSecretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		nonce, err := utils.GenerateRandomNonce()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": " bad nonce generation"})
			return
		}

		userRepo.UpdateUserNonce(retrievedUser.ID, nonce)

		c.JSON(http.StatusOK, gin.H{"token": tokenString})

	} else {

		fmt.Println("Signature verification failed")

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Signature verification failed"})
	}
}
