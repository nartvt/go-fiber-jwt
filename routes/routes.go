package routes

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-jwt/controller"
)

func Setup(app *fiber.App) {
	app.Get("/", controller.Hello)

	app.Post("/api/register", controller.Register)
	app.Post("/api/login", controller.Login)

	app.Get("/api/user", controller.User)

}
