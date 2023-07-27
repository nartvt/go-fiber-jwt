package controller

import (
	"fmt"
	"go-fiber-jwt/app/middleware"
	"go-fiber-jwt/app/request"
	"go-fiber-jwt/app/response"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Register(ctx *fiber.Ctx) error {
	fmt.Println("Fiber Register")
	var data request.UserRequest
	if err := data.Bind(ctx); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	err := middleware.RegisterUser(data, ctx)
	if err != nil {
		return middleware.ResponseWithStatus(ctx, err.Code, err.Message)
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
}

func Login(ctx *fiber.Ctx) error {
	fmt.Println("Fiber LOGIN")
	var data request.UserRequest
	if err := data.Bind(ctx); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err := middleware.Authenticated(data, ctx)
	if err != nil {
		return middleware.ResponseWithStatus(ctx, err.Code, err.Message)
	}

	return response.ResponseUserOK(ctx)
}

func User(ctx *fiber.Ctx) error {
	user, err := middleware.IsUser(ctx)
	if err != nil {
		return middleware.ResponseWithStatus(ctx, err.Code, err.Message)
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success", "user": user})
}

func Logout(ctx *fiber.Ctx) error {
	fmt.Println("Fiber LOGOUT")
	_, err := middleware.IsUser(ctx)
	if err != nil {
		return middleware.ResponseWithStatus(ctx, err.Code, err.Message)
	}
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
