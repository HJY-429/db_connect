package main

import (
	"fmt"
	"os"
	"strings"

	"tidb-gin-demo/config"
	"tidb-gin-demo/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	// initialize database
	config.InitDB()

	// set gin mode from environment (default: debug)
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}

	// create gin router without pre-attached middleware
	r := gin.New()
	// attach Logger and Recovery exactly once
	r.Use(gin.Logger(), gin.Recovery())

	// configure trusted proxies (comma-separated list in TRUSTED_PROXIES)
	var trusted []string
	if tp := os.Getenv("TRUSTED_PROXIES"); tp != "" {
		for _, p := range strings.Split(tp, ",") {
			trusted = append(trusted, strings.TrimSpace(p))
		}
	} else {
		// sensible default: localhost only
		trusted = []string{"127.0.0.1"}
	}
	if err := r.SetTrustedProxies(trusted); err != nil {
		fmt.Println("Warning: failed to set trusted proxies:", err)
	}

	// create controller instances
	userController := &controllers.UserController{}

	// user routes
	userRoutes := r.Group("/api/users")
	{
		userRoutes.POST("/", userController.CreateUser)      // create user
		userRoutes.GET("/", userController.GetUsers)         // get all users
		userRoutes.GET("/:id", userController.GetUser)       // get user by id
		userRoutes.PUT("/:id", userController.UpdateUser)    // update user by id
		userRoutes.DELETE("/:id", userController.DeleteUser) // delete user by id
	}

	// health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "TiDB Gin Demo is running!",
		})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// run the server
	r.Run(":8080")
}
