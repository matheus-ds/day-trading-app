package main

import (
	"day-trading-app/backend/cmd/server"
	"day-trading-app/backend/pkg/logger"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Error("No .env file found")
	}
}

func main() {
	logger.Info("Server started")
	server.Run()
	logger.Info("Server stopped")
}
