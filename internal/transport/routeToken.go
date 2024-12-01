package transport

import (
	"mymod/internal/services"

	"github.com/gofiber/fiber/v2"
)

// выдаёт аксес и рефреш токен по guid
func getToken(c *fiber.Ctx) error {

	var GUID = c.Query("GUID", "")
	if GUID == "" {
		return c.Status(400).JSON("empty query")
	}

	res, err := services.RouteGetToken(GUID)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.Status(200).JSON(res)
}

// делает рефреш токена. выдаёт новый аксес токен, проверяет совпадения и айпи
func refreshToken(c *fiber.Ctx) error {

	var refreshToken = c.Query("refresh", "")
	var accessToken = c.Query("access", "")
	if refreshToken == "" || accessToken == "" {
		return c.Status(400).JSON("empty query")
	}

	res, err := services.RouteRefreshToken(accessToken, refreshToken)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.Status(200).JSON(res)
}
