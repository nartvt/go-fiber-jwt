package routes

import (
	"go-fiber-jwt/app/controller"
	"go-fiber-jwt/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Use(middleware.CorsFilter())

	api := app.Group("/api")
	POST(api, "/register", controller.Register)

	POST(api, "/login", controller.Login)

	GET(api, "/user", controller.User)

	PUT(api, "/logout", controller.Logout)
}
