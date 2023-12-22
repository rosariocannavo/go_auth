package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rosariocannavo/go_auth/internal/nats"
)

func RoleAuth(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.MustGet("claims").(jwt.MapClaims)

		if !ok {

			message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "middleware/RoleAuth", http.StatusForbidden, "error: Invalid token claims")
			nats.NatsConnection.PublishMessage(message)

			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Extract the role from the token claims
		tokenRole, ok := claims["role"].(string)
		if !ok {

			message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "middleware/RoleAuth", http.StatusForbidden, "error: Role information not found")
			nats.NatsConnection.PublishMessage(message)

			c.JSON(http.StatusForbidden, gin.H{"error": "Role information not found"})
			c.Abort()
			return
		}

		// Authorize based on the extracted role
		if tokenRole != role {

			message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "middleware/RoleAuth", http.StatusForbidden, "error: User not authorized to access this resource")
			nats.NatsConnection.PublishMessage(message)

			c.JSON(http.StatusForbidden, gin.H{"error": "User not authorized to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
