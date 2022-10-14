package routes

import (
	"github.com/Ankan002/compiler-api/controllers/compiler"
	"github.com/gofiber/fiber/v2"
)

func CompilerRouter(router fiber.Router) {
	router.Post("/compiler", compiler.Compiler)
}
