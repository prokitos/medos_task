package transport

import (
	"errors"
	"mymod/internal/services"

	"github.com/gofiber/fiber/v2"
)

// выдаёт аксес и рефреш токен по guid
func getToken(c *fiber.Ctx) error {

	var GUID = c.Query("GUID", "")

	services.RouteGetToken(GUID)

	return errors.New(GUID)
}

// делает рефреш токена. выдаёт новый аксес токен, проверяет совпадения и айпи
func refreshToken(c *fiber.Ctx) error {

	var refreshToken = c.Query("refresh", "")
	var accessToken = c.Query("access", "")

	services.RouteRefreshToken(accessToken, refreshToken)

	return errors.New(accessToken + refreshToken)
}
