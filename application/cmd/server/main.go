// main.go

// @title Ihale API
// @description This is the Ihale API server.
// @version 1
// @BasePath /api/v1
package main

import (
	"blockchain-smart-tender-platform/internal/server"
	"blockchain-smart-tender-platform/pkg/logger"
	_ "blockchain-smart-tender-platform/platform/channel"

	_ "blockchain-smart-tender-platform/cmd/server/docs"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger.Init()
	server.Run()
}
