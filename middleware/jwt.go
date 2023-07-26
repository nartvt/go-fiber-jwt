package middleware

import (
	"encoding/json"
	"log"
	"time"

	"go-fiber-jwt/app/models"
	"go-fiber-jwt/app/repo"
	"go-fiber-jwt/app/request"
	"go-fiber-jwt/infra"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	SecretKey     = "secret"
	jwtExpireTime = time.Hour * 2
)

type Error *fiber.Error
type Claims jwt.MapClaims
type Token jwt.Token

func RegisterUser(input request.UserRequest, ctx *fiber.Ctx) Error {
	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), 4)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if len(input.Email) <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "email invalid")
	}

	user, uerr := repo.UserRepository.GetUserByEmail(input.Email)
	if uerr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if user != nil {
		return fiber.NewError(fiber.StatusForbidden, "this user is exists")
	}

	user = &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: password,
	}
	err = infra.DB.Create(user).Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	token, err := generateJWTToken(user.Id, user.Email)
	if err != nil {
		log.Printf("token.SignedString: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(jwtExpireTime),
		HTTPOnly: true,
	}
	ctx.Cookie(&cookie)
	return nil
}

func Authenticated(data request.UserRequest, ctx *fiber.Ctx) Error {
	isUser, err := ValidateLogined(ctx)
	if err == nil && isUser {
		return nil
	}

	user, uerr := repo.UserRepository.GetUserByEmail(data.Email)
	if uerr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, uerr.Error())
	}

	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "user invalid")
	}

	if err := verifyPassword(user.Password, data.Password); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	token, uerr := generateJWTToken(user.Id, user.Email)
	if uerr != nil {
		log.Printf("token.SignedString: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, uerr.Error())
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(jwtExpireTime),
		HTTPOnly: true,
	}
	ctx.Cookie(&cookie)
	userData, uerr := json.Marshal(user)
	if uerr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, uerr.Error())
	}
	ctx.Set("user", string(userData))
	return nil
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func verifyPassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

func ValidateLogined(ctx *fiber.Ctx) (bool, Error) {
	return IsUser(ctx)
}

func IsUser(ctx *fiber.Ctx) (bool, Error) {
	if ctx == nil {
		return false, fiber.NewError(fiber.StatusBadRequest, "bad request")
	}

	claims, err := validateJWT(ctx)
	if err != nil {
		return false, err
	}
	if len(claims) <= 0 {
		return false, nil
	}
	ctx.Set("claims", claims["jwt"].(string))
	return true, nil
}

func validateJWT(ctx *fiber.Ctx) (Claims, Error) {
	claims, token, err := destructJWTToken(ctx)
	if err != nil {
		return Claims{}, err
	}
	if !token.Valid {
		return Claims{}, fiber.NewError(fiber.StatusUnauthorized, "token invalid")
	}

	expireDate, uerr := claims.GetExpirationTime()
	if uerr != nil {
		return Claims{}, fiber.NewError(fiber.StatusUnauthorized, "user expire")
	}

	if time.Now().Sub(expireDate.Time) >= jwtExpireTime {
		return Claims{}, fiber.NewError(fiber.StatusUnauthorized, "user expire")
	}
	return Claims{}, nil
}

func destructJWTToken(ctx *fiber.Ctx) (jwt.MapClaims, *jwt.Token, Error) {
	cookie := ctx.Cookies("jwt")
	if len(cookie) <= 0 {
		return jwt.MapClaims{}, &jwt.Token{}, fiber.NewError(fiber.StatusUnauthorized, "jwt token invalid")
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return jwt.MapClaims{}, &jwt.Token{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return claims, token, nil
}

func generateJWTToken(userId int, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"exp":     time.Now().Add(jwtExpireTime).Unix(), // Token expires in 24 hours
	})
	return token.SignedString([]byte(SecretKey))
}

func ResponseWithStatus(ctx *fiber.Ctx, code int, message string) error {
	return ctx.Status(code).JSON(fiber.Map{"message": message})
}
