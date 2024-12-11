package main

import (
	"airbnb-api/controllers"
	"airbnb-api/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	database.ConnectDB()

	// Initialize Gin router
	r := gin.Default()

	// Routes
	r.GET("/:room_id", controllers.GetRoomMetrics)

	// Start server
	r.Run(":8080")
}
