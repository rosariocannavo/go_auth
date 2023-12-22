package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rosariocannavo/go_auth/config"
	"github.com/rosariocannavo/go_auth/internal/db"
	"github.com/rosariocannavo/go_auth/internal/models"
	"github.com/rosariocannavo/go_auth/internal/nats"
	"github.com/rosariocannavo/go_auth/internal/repositories"
	"github.com/rosariocannavo/go_auth/internal/utils"
)

func HandleRegistration(c *gin.Context) {
	userRepo := repositories.NewUserRepository(db.Client)

	var userForm models.UserForm

	//retrieve the partial user information from form
	if err := c.BindJSON(&userForm); err != nil {

		message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "registration_handler/HandleRegistration", http.StatusBadRequest, "error: Invalid request payload")
		nats.NatsConnection.PublishMessage(message)

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	//check if the user is already registered
	isPresent, err := userRepo.CheckIfUserIsPresent(userForm.Username, userForm.MetamaskAddress)

	if err != nil {

		message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "registration_handler/HandleRegistration", http.StatusInternalServerError, "error: Error fetching database")
		nats.NatsConnection.PublishMessage(message)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching database"})
		return
	}

	if isPresent {

		message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "registration_handler/HandleRegistration", http.StatusForbidden, "error: User already present")
		nats.NatsConnection.PublishMessage(message)

		c.JSON(http.StatusForbidden, gin.H{"error": "User already present"})
		return

	} else {

		//if user is not present
		//hash his psw and store him in the db
		//give him a role based on # of transaction
		//give him a nonce

		var user models.User

		// hash the user password
		hashedPwd, err := utils.HashPassword(userForm.Password)

		if err != nil {

			message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "registration_handler/HandleRegistration", http.StatusInternalServerError, "error: Error hashing password")
			nats.NatsConnection.PublishMessage(message)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		// generate nonce for metamask sign auth
		nonce, err := utils.GenerateRandomNonce()

		if err != nil {

			message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "registration_handler/HandleRegistration", http.StatusInternalServerError, "error: Bad nonce generation")
			nats.NatsConnection.PublishMessage(message)

			c.JSON(http.StatusInternalServerError, gin.H{"error": " Bad nonce generation"})
			return
		}

		user.Username = userForm.Username
		user.Password = hashedPwd
		user.MetamaskAddress = userForm.MetamaskAddress
		user.Nonce = nonce

		//TODO
		// generate role based on nonce

		//conta il numero di transazioni fatte
		//payload := strings.NewReader(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["%s", "latest"],"id":1}`, userForm.MetamaskAddress))
		payload := strings.NewReader(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBalance","params":["%s", "latest"],"id":1}`, userForm.MetamaskAddress))

		// Sending the HTTP POST request to the Ganache endpoint
		resp, err := http.Post(config.GanacheURL, "application/json", payload)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// Reading the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("BODY", string(body))

		//per adesso vedo banalmente se l'account corrisponde a quello di ganache che ho su metamask
		if strings.EqualFold(user.MetamaskAddress, "0x58ad8fEA5d85EDD13C05dC116198801Ff53679B2") {
			user.Role = models.Admin
		} else {
			user.Role = models.NormalUser
		}

		// write the user in the database
		userRepo.CreateUser(&user)

		message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "registration_handler/HandleRegistration", http.StatusOK, "message: User registered succesfully")
		nats.NatsConnection.PublishMessage(message)

		c.JSON(http.StatusOK, gin.H{"message": "User registered succesfully"})
	}
}
