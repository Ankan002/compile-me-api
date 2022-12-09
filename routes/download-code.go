package routes

import (
	download_code "github.com/Ankan002/compiler-api/controllers/download-code"
	"github.com/gofiber/fiber/v2"
)

func DownloadCode(router fiber.Router) {
	router.Post("/download-code", download_code.DownloadCodeFile)
}
