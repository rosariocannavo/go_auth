package main

import (
	"fmt"
	"log"

	"github.com/rosariocannavo/go_auth/internal/db"
	"github.com/rosariocannavo/go_auth/internal/handlers"
	"github.com/rosariocannavo/go_auth/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseDB()

	fmt.Println("Connected to MongoDB!")

	// Apply middleware to all routes

	// Login endpoint
	r.POST("/login", handlers.HandleLogin)

	//registration endpoint
	r.POST("/registration", handlers.HandleRegistration)

	//Protected middleware User endpoints
	userRoutes := r.Group("/users")
	userRoutes.Use(middleware.Authenticate())
	userRoutes.Use(middleware.RoleAuth("user"))
	{
		userRoutes.GET("/data" /*middleware.Authenticate(),*/, handlers.GetUserData)
	}

	//Protected middleware Admin endpoints
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.Authenticate())
	userRoutes.Use(middleware.RoleAuth("admin"))
	{
		adminRoutes.GET("/data" /*middleware.Authenticate(),*/, handlers.GetAdminData)
	}

	// Run the server
	_ = r.Run(":8080")

}
