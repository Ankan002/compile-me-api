package routes

import (
	"github.com/Ankan002/compiler-api/controllers/languages"
	"github.com/gofiber/fiber/v2"
)

func LanguagesRouter(router fiber.Router) {
	router.Get("/languages", languages.GetLanguages)
}
