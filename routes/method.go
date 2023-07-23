package routes

import "github.com/gofiber/fiber/v2"

func GET(r fiber.Router, path string, handler fiber.Handler) {
	r.Get(path, handler)
}

func POST(r fiber.Router, path string, handler fiber.Handler) {
	r.Post(path, handler)
}

func PUT(r fiber.Router, path string, handler fiber.Handler) {
	r.Post(path, handler)
}
