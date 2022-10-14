package main

import (
	"github.com/Ankan002/compiler-api/config"
	"github.com/Ankan002/compiler-api/routes"
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

	router := app.Group("/api")

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"success": true,
			"message": "Welcome to the compiler API",
		})
	})

	routes.CompilerRouter(router)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
