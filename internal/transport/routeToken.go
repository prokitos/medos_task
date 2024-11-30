package transport

import (
	"mymod/internal/services"

	"github.com/gofiber/fiber/v2"
)

// выдаёт аксес и рефреш токен по guid
func getToken(c *fiber.Ctx) error {

	var GUID = c.Query("GUID", "")
	res, err := services.RouteGetToken(GUID)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(res)
}

// делает рефреш токена. выдаёт новый аксес токен, проверяет совпадения и айпи
func refreshToken(c *fiber.Ctx) error {

	var refreshToken = c.Query("refresh", "")
	var accessToken = c.Query("access", "")

	res, err := services.RouteRefreshToken(accessToken, refreshToken)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(res)
}
