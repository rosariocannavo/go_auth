package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rosariocannavo/go_auth/internal/nats"
)

func GetCookieHandler(c *gin.Context) {
	jwtCookie, err := c.Cookie("jwtToken")
	if err != nil {

		message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d", time.Now().UTC().Format(time.RFC3339), "cookie_handler/GetCookieHandler", http.StatusNotFound)
		nats.NatsConnection.PublishMessage(message)

		c.String(http.StatusNotFound, "jwt Cookie not found")
		return
	}

	accountAddressCookie, err := c.Cookie("accountAddress")
	if err != nil {

		message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d", time.Now().UTC().Format(time.RFC3339), "cookie_handler/GetCookieHandler", http.StatusNotFound)
		nats.NatsConnection.PublishMessage(message)

		c.String(http.StatusNotFound, "account Cookie not found")
		return
	}

	message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "cookie_handler/GetCookieHandler", http.StatusOK, "token: jwtCookie, account: accountAddressCookie")
	nats.NatsConnection.PublishMessage(message)

	c.JSON(http.StatusOK, gin.H{"token": jwtCookie, "account": accountAddressCookie})
}
