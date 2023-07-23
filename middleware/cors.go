package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CorsFilter() fiber.Handler {
	return cors.New(corsConfig())
}

const appName = "api gateway version 1.0"

func corsConfig() cors.Config {
	return cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,HEAD,DELETE,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}
}

func AppName() fiber.Config {
	return fiber.Config{
		BodyLimit:   2 * 1024 * 1024,
		Concurrency: 100,
		AppName:     appName,
	}
}
