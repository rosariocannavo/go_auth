package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAdminData(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username not found in context"})
		return
	}
	// Access the username from the context (authenticated user)
	// For demonstration purposes, just returning the username in response
	c.JSON(http.StatusOK, gin.H{"message": "Data fetched successfully", "username": username})
}
