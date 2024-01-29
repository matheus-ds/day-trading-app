package main

import (
	"github.com/joho/godotenv"
	"github.com/matheus-ds/day-trading-app/backend/cmd/server"
	"github.com/matheus-ds/day-trading-app/backend/pkg/logger"
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
