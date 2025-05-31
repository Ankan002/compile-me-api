package main

import (
	"log"
	"os"

	"github.com/Ankan002/compiler-api/config"
	f_lambda "github.com/Ankan002/compiler-api/lambda"
	"github.com/Ankan002/compiler-api/routes"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func serverfulHandler() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	router := app.Group("/api")

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"success": true,
			"message": "Welcome to the Compile Me API",
		})
	})

	routes.CompilerRouter(router)
	routes.LanguagesRouter(router)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}

func main() {
	if os.Getenv("GO_ENV") != "production" {
		config.LoadEnv()
	}

	if os.Getenv("INVOCATION") == "non-lambda" {
		serverfulHandler()
	}

	if os.Getenv("INVOCATION") == "lambda" {
		lambda.Start(f_lambda.CompilationLambda)
	}
}
