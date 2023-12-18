package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rosariocannavo/go_auth/internal/db"
	"github.com/rosariocannavo/go_auth/internal/handlers"
	"github.com/rosariocannavo/go_auth/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	// Load Static files
	r.LoadHTMLGlob("../../templates/*.html")
	r.Static("/css", "../../templates/css")
	r.Static("/js", "../../templates/js")

	// Connect to DB
	err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseDB()

	fmt.Println("Connected to MongoDB!")

	r.Use(middleware.NewRateLimitMiddleware().Handler())

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", gin.H{})
	})

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})

	// Cookie endpoint to retrieve account and jwt
	r.GET("/get-cookie", handlers.GetCookieHandler)

	// Login endpoint
	r.POST("/login", handlers.HandleLogin)

	//registration endpoint
	r.POST("/registration", handlers.HandleRegistration)

	//metamask signature verification endpoint
	r.POST("/verify-signature", handlers.HandleverifySignature)

	//all this logic must be in the other service and protected by proxy
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

	//TODO: this route use the proxy + cb launched by button in login
	r.GET("/app/*proxyPath", middleware.Authenticate(), handlers.ProxyHandler) //handler of the proxyy

	// Run the server
	_ = r.Run(":8080")

}
