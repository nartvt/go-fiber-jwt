package request

import "github.com/gofiber/fiber/v2"

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (input *UserRequest) Bind(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(input); err != nil {
		return err
	}
	return nil
}
