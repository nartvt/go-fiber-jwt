package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go-fiber-jwt/database"
	"go-fiber-jwt/models"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

func Hello(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello, World !")
}

func Register(ctx *fiber.Ctx) error {
	var data map[string]string

	err := ctx.BodyParser(&data)
	if err != nil {
		return err
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 4)
	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	database.DB.Create(&user)

	return ctx.JSON(user)
}

const SecretKey = "secret"

func Login(ctx *fiber.Ctx) error {
	var data map[string]string
	if err := ctx.BodyParser(&data); err != nil {
		return err
	}
	var user models.User

	database.DB.Where("email=?", data["email"]).First(&user)
	if user.Id == 0 {
		ctx.Status(fiber.StatusFound)
		return ctx.JSON(fiber.Map{"message": "User not found"})
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{"message": "Password incorrect!"})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // one day
	})

	token, err1 := claims.SignedString([]byte(SecretKey))

	if err1 != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"message": "Could not login"})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	ctx.Cookie(&cookie)
	return ctx.JSON(fiber.Map{"message": "success"})
}

func User(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		ctx.Status(fiber.StatusUnauthorized)
		return ctx.JSON(fiber.Map{"message": "unauthenticate"})
	}
	claims := token.Claims
	return ctx.JSON(claims)
}
