package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCookieHandler(c *gin.Context) {
	jwtCookie, err := c.Cookie("jwtToken")
	if err != nil {
		c.String(http.StatusNotFound, "jwt Cookie not found")
		return
	}

	accountAddressCookie, err := c.Cookie("accountAddress")
	if err != nil {
		c.String(http.StatusNotFound, "account Cookie not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": jwtCookie, "account": accountAddressCookie})
}
