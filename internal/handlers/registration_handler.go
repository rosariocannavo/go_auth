package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rosariocannavo/go_auth/config"
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

	fmt.Println("ADDRESS FROM REGI ", userForm.MetamaskAddress)
	//check if the user is already registered
	isPresent, err := userRepo.CheckIfUserIsPresent(userForm.Username, userForm.MetamaskAddress)

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
		//give him a role based on # of transaction
		//give him a nonce

		var user models.User

		// hash the user password
		hashedPwd, err := utils.HashPassword(userForm.Password)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		// generate nonce for metamask sign auth
		nonce, err := utils.GenerateRandomNonce()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": " bad nonce generation"})
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
			fmt.Println("admin logger")
			user.Role = models.Admin
		} else {
			fmt.Println("user logger")
			user.Role = models.NormalUser
		}

		// write the user in the database
		userRepo.CreateUser(&user)

		c.JSON(http.StatusOK, gin.H{"message": "user registered succesfully"})
	}
}
