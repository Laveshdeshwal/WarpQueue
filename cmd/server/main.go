package main

import (
	"WarpQueue/internal/api"
	"WarpQueue/internal/logger"
	"WarpQueue/internal/queue"
	"log"
	"net/http"

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

	q := queue.NewMemoryQueue()
	server := api.NewServer(q)
	logger.InitLogger("app-log")

	logger.Info("Server running on :8080")

	err := http.ListenAndServe(":"+viper.GetString("SERVER_PORT"), server.Routes())
	if err != nil {
		log.Fatal(err)
	}

	//app.Run(":8080")
}
