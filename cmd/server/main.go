package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// Load Environment Variables
	envPath := ".env"

	viper.SetConfigFile(envPath)
	viper.ReadInConfig()
	viper.AutomaticEnv()

	app := gin.Default()

	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Server": "go-warp-queue-server",
			"Status": "running",
		})
	})

	app.Run(":8080")
}
