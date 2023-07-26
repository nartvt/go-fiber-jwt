package response

import (
	"encoding/json"
	"go-fiber-jwt/app/models"

	"github.com/gofiber/fiber/v2"
)

type UserResponse struct {
	Id    int
	Name  string
	Token string
}

func ResponseUserOK(ctx *fiber.Ctx) error {
	user := ctx.Get("user")
	if len(user) <= 0 {
		return nil
	}
	var userModel models.User
	err := json.Unmarshal([]byte(user), &userModel)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "user invalid")
	}
	userResp := UserResponse{
		Id:    userModel.Id,
		Name:  userModel.Name,
		Token: ctx.Cookies("jwt"),
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success", "user": userResp})
}
