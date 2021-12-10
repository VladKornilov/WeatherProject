package main

import (
	"../../internal/server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	router := gin.Default()

	server.StartServer(router)
}
