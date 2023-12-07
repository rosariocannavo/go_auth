package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	auth "github.com/rosariocannavo/go_auth/config"
	"github.com/rosariocannavo/go_auth/internal/db"
	"github.com/rosariocannavo/go_auth/internal/models"
	"github.com/rosariocannavo/go_auth/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

func checkSig(from, sigHex string, msg []byte) bool {
	sig := hexutil.MustDecode(sigHex)

	msg = accounts.TextHash(msg)
	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}

	fmt.Printf("ECDSA Signature: %x\n", sig)
	fmt.Printf("  R: %x\n", sig[0:32])  // 32 bytes
	fmt.Printf("  S: %x\n", sig[32:64]) // 32 bytes
	fmt.Printf("  V: %x\n", sig[64:])

	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	fmt.Println("recovered", recoveredAddr.Hex())
	return strings.EqualFold(from, recoveredAddr.Hex())
}

func generateRandomNonce2() (string, error) {
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

	return nonceString, nil
}

var retrievedUser models.User

func HandleLogin(c *gin.Context) {
	userRepo := repositories.NewUserRepository(db.Client)
	var userForm models.UserForm

	if err := c.BindJSON(&userForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	retrievedUser, err := userRepo.FindUser(userForm.Username)

	fmt.Println("username", userForm.Username)
	if err != nil {
		fmt.Println("USER NOT PRESENT")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not present"})
		return
	}

	errf := bcrypt.CompareHashAndPassword([]byte(retrievedUser.Password), []byte(userForm.Password))

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

	fmt.Println("mesg " + requestData.Nonce)
	fmt.Println("addr " + requestData.Address)
	fmt.Println("sig " + requestData.Signature + "\n")

	signatureVerified := checkSig(requestData.Address, requestData.Signature, []byte(requestData.Nonce))

	if signatureVerified {
		fmt.Println("Signature verification success")

		// Generate JWT token upon successful authentication
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":              retrievedUser.ID,
			"username":        retrievedUser.Username,
			"metamaskAddress": retrievedUser.MetamaskAddress,
			"nonce":           retrievedUser.Nonce,
			"role":            retrievedUser.Role,

			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString(auth.SecretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		nonce, err := generateRandomNonce2()
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
