package routes

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-jwt/controller"
)

func Setup(app *fiber.App) {

	api := app.Group("/api")

	api.Post("/register", controller.Register)

	api.Post("/login", controller.Login)

	api.Get("/user", controller.User)

	api.Post("/logout", controller.Logout)

}
