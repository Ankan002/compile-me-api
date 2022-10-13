package main

import (
	"github.com/Ankan002/compiler-api/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
)

func main() {
	if os.Getenv("GO_ENV") != "production" {
		config.LoadEnv()
	}

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"success": true,
			"message": "Welcome to the compiler API",
		})
	})

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
