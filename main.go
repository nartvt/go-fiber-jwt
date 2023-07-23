package main

import (
	"go-fiber-jwt/infra"
	"go-fiber-jwt/middleware"
	"go-fiber-jwt/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New(middleware.AppName())
	routes.Setup(app)

	setupDatabase()
	defer infra.CloseConnect()

	app.Listen(":3000")
}

func setupDatabase() {
	infra.Connect()
}
