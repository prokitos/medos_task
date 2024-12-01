package transport

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// проверка что роуты доступны.
func TestHttpConnect(t *testing.T) {

	tests := []struct {
		description  string
		route        string
		expectedCode int
		method       string
	}{
		{
			description:  "get HTTP status 404, when route is not exists",
			route:        "/not-found",
			expectedCode: 404,
			method:       "GET",
		},
		{
			description:  "get HTTP status 400, send empty query to route",
			route:        "/getToken",
			expectedCode: 400,
			method:       "GET",
		},
		{
			description:  "get HTTP status 400, send empty query to route",
			route:        "/refreshToken",
			expectedCode: 400,
			method:       "GET",
		},
	}

	app := fiber.New()
	SetHandlers(app)

	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.route, nil)
		resp, _ := app.Test(req, 1)
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
