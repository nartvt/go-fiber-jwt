package controller

import (
	"fmt"
	"go-fiber-jwt/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Register(ctx *fiber.Ctx) error {
	fmt.Println("Fiber Register")
	err := middleware.RegisterUser(ctx)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
}

func Login(ctx *fiber.Ctx) error {
	fmt.Println("Fiber LOGIN")
	var data map[string]string
	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	err := middleware.Authenticated(data, ctx)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}

func User(ctx *fiber.Ctx) error {
	fmt.Println("Fiber Get user info success")
	user, err := middleware.IsUser(ctx)
	fmt.Println("39 - ", err != nil)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success", "user": user})
}

func Logout(ctx *fiber.Ctx) error {
	fmt.Println("Fiber logout")
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	ctx.Cookie(&cookie)
	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}
