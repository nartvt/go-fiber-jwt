package main

import (
	"fmt"
	"go-fiber-jwt/config"
	"go-fiber-jwt/infra"
	"go-fiber-jwt/internal/routes"
	"go-fiber-jwt/middleware"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.GetConfig()

	app := fiber.New(middleware.AppName())
	routes.Setup(app)

	setupDatabase()
	defer infra.CloseConnect()

	app.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
}

func setupDatabase() {
	infra.Connect()
}
