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
	"github.com/rosariocannavo/go_auth/internal/nats"
	"github.com/rosariocannavo/go_auth/internal/repositories"

	"github.com/rosariocannavo/go_auth/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var retrievedUser *models.User

func HandleLogin(c *gin.Context) {
	userRepo := repositories.NewUserRepository(db.Client)
	var userForm models.UserForm

	if err := c.BindJSON(&userForm); err != nil {

		message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "login_handler/HandleLogin", http.StatusBadRequest, "error: Invalid request payload")
		nats.NatsConnection.PublishMessage(message)

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var errdb error
	retrievedUser, errdb = userRepo.FindUser(userForm.Username)

	fmt.Println("username", userForm.Username)
	if errdb != nil {

		message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "login_handler/HandleLogin", http.StatusBadRequest, "error: User not present")
		nats.NatsConnection.PublishMessage(message)

		c.JSON(http.StatusBadRequest, gin.H{"error": "User not present"})
		return
	}

	errf := utils.CompareHashPwd(retrievedUser.Password, userForm.Password)

	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!

		message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "login_handler/HandleLogin", http.StatusUnauthorized, "error: Invalid password")
		nats.NatsConnection.PublishMessage(message)

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: {%s:%s}", time.Now().UTC().Format(time.RFC3339), "login_handler/HandleLogin", http.StatusAccepted, "Nonce", retrievedUser.Nonce)
	nats.NatsConnection.PublishMessage(message)

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

		message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "login_handler/HandleverifySignature", http.StatusBadRequest, "error: Invalid JSON")
		nats.NatsConnection.PublishMessage(message)

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

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

			message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "login_handler/HandleverifySignature", http.StatusInternalServerError, "error: Failed to generate token")
			nats.NatsConnection.PublishMessage(message)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		nonce, err := utils.GenerateRandomNonce()
		if err != nil {

			message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "login_handler/HandleverifySignature", http.StatusInternalServerError, "error: Bad nonce generation")
			nats.NatsConnection.PublishMessage(message)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Bad nonce generation"})
			return
		}

		userRepo.UpdateUserNonce(retrievedUser.ID, nonce)

		//set cookie for the session and return token to client
		jwtCookie := &http.Cookie{
			Name:  "jwtToken",
			Value: tokenString,
			Path:  "/",
		}
		http.SetCookie(c.Writer, jwtCookie)

		accountCookie := &http.Cookie{
			Name:  "accountAddress",
			Value: retrievedUser.MetamaskAddress,
			Path:  "/",
		}
		http.SetCookie(c.Writer, accountCookie)

		message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: {%s:%s, %s:%s}", time.Now().UTC().Format(time.RFC3339), "login_handler/HandleverifySignature", http.StatusOK, "token", tokenString, "role", retrievedUser.Role)
		nats.NatsConnection.PublishMessage(message)

		c.JSON(http.StatusOK, gin.H{"token": tokenString, "role": retrievedUser.Role})

	} else {

		message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "login_handler/HandleverifySignature", http.StatusUnauthorized, "error: Signature verification failed")
		nats.NatsConnection.PublishMessage(message)

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Signature verification failed"})
	}
}
