package transport

import "github.com/gofiber/fiber/v2"

func SetHandlers(instance *fiber.App) {

	instance.Get("/getToken", getToken)
	instance.Get("/refreshToken", refreshToken)

}
