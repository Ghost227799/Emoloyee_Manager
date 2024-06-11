package main

import (
	config "Shaunak/Employee_manager/Config"
	handlers "Shaunak/Employee_manager/Services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// Initialize Gin router
	r := gin.Default()
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	// Connect to the database
	err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer config.CloseDB()
	// Middleware to inject the database connection into the context
	r.Use(func(c *gin.Context) {
		c.Set("db", config.DB())
		c.Next()
	})

	// Register employee handlers
	handlers.RegisterEmployeeHandlers(r)

	// Run server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
