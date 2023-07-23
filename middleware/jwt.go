package middleware

import (
	"errors"
	"log"
	"time"

	"go-fiber-jwt/app/models"
	"go-fiber-jwt/infra"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	SecretKey = "secret"
	expire    = time.Hour * 2
)

func RegisterUser(ctx *fiber.Ctx) error {
	var data map[string]string

	err := ctx.BodyParser(&data)
	if err != nil {
		return err
	}
	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 4)
	if err != nil {
		return err
	}
	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}
	err = infra.DB.Create(&user).Error
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	token, err := generateJWTToken(user.Id, user.Email)
	if err != nil {
		log.Printf("token.SignedString: %v", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	ctx.Cookie(&cookie)
	return nil
}

func Authenticated(data map[string]string, ctx *fiber.Ctx) error {
	var user models.User
	email := data["email"]
	err := infra.DB.First(&user).Where("email = ?", email).Error
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if err := verifyPassword(user.Password, data["password"]); err != nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	token, err := generateJWTToken(user.Id, user.Email)
	if err != nil {
		log.Printf("token.SignedString: %v", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	ctx.Cookie(&cookie)
	return nil
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func verifyPassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

func IsUser(ctx *fiber.Ctx) (interface{}, error) {
	if ctx == nil {
		return nil, errors.New("bad request")
	}
	user, err := validateJWT(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthenticate"})
	}
	return user, nil
}

func validateJWT(ctx *fiber.Ctx) (*models.User, error) {
	cookie := ctx.Cookies("jwt")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
	}
	var user *models.User
	userId := claims["user_id"]
	err = infra.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	if user == nil {
		return nil, ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthenticate"})
	}
	email := claims["email"]
	if email != user.Email {
		return nil, ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthenticate"})
	}
	return user, nil
}

func generateJWTToken(userId int, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"exp":     time.Now().Add(expire).Unix(), // Token expires in 24 hours
	})
	return token.SignedString([]byte(SecretKey))
}
