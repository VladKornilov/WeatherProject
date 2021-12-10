package main

import (
	"../../internal/server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	app, err := server.CreateApplication()

	if err != nil {
		panic(err)
	}
	defer func(app *server.Application) {
	}(app)

	app.StartServer()
}
