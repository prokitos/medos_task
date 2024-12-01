package services

import (
	"errors"
	"mymod/internal/models"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// проверка аксес токена
func TestCheckAccessToken(t *testing.T) {

	tests := []struct {
		description string
		expected    error
		token       string
	}{
		{
			description: "normal token check",
			expected:    nil,
			token:       createTokenAccess("8e2b4136-829c-11eb-8dcd-0242ac130003", "127.0.0.1"),
		},
		{
			description: "old token check",
			expected:    errors.New("token is expired by 0s"),
			token:       oldTokenCreate(),
		},
		{
			description: "error token check",
			expected:    errors.New("token contains an invalid number of segments"),
			token:       "zxczg",
		},
	}

	for _, test := range tests {

		resp := checkAccessToken(test.token)
		assert.Equalf(t, test.expected, resp, test.description)
	}
}

// проверка рефреш токена
func TestCheckRefreshToken(t *testing.T) {

	goodAccessToken := createTokenAccess("8e2b4136-829c-11eb-8dcd-0242ac130003", "127.0.0.1")

	tests := []struct {
		description string
		expected    error
		refToken    string
		accToken    string
	}{
		{
			description: "normal token check",
			expected:    nil,
			refToken:    createTokenRefresh("8e2b4136-829c-11eb-8dcd-0242ac130003", goodAccessToken),
			accToken:    goodAccessToken,
		},
		{
			description: "different access token check",
			expected:    errors.New("access token missmatch"),
			refToken:    createTokenRefresh("8e2b4136-829c-11eb-8dcd-0242ac130003", createTokenAccess("8e2b4136-829c-11eb-8dcd-0242ac130003", "127.0.0.1")),
			accToken:    createTokenAccess("8e2b4136-829c-11eb-8dcd-0242ac130003", "127.0.0.2"),
		},
		{
			description: "error token check",
			expected:    errors.New("token contains an invalid number of segments"),
			refToken:    "zxczg",
			accToken:    "dsgs",
		},
	}

	for _, test := range tests {

		resp := checkRefreshTokens(test.refToken, test.accToken)
		assert.Equalf(t, test.expected, resp, test.description)
	}
}

// проверка текущего айпи и айпи токена
func TestGetDataFromAccess(t *testing.T) {

	tests := []struct {
		description string
		expected    error
		token       string
	}{
		{
			description: "check curr ip",
			expected:    nil,
			token:       createTokenRefresh("8e2b4136-829c-11eb-8dcd-0242ac130003", createTokenAccess("8e2b4136-829c-11eb-8dcd-0242ac130003", "127.0.0.1")),
		},
	}

	for _, test := range tests {

		resp := checkTokenIp(test.token, "127.0.0.1")
		assert.Equalf(t, test.expected, resp, test.description)
	}
}

// получить guid из аксес токена
func TestGetGuidByToken(t *testing.T) {

	tests := []struct {
		description string
		expected    error
		token       string
	}{
		{
			description: "get guid",
			expected:    nil,
			token:       createTokenAccess("8e2b4136-829c-11eb-8dcd-0242ac130003", "127.0.0.1"),
		},
	}

	for _, test := range tests {

		data, resp := getGuid(test.token)
		assert.Equalf(t, test.expected, resp, test.description)
		assert.Equalf(t, "8e2b4136-829c-11eb-8dcd-0242ac130003", data, test.description)

	}
}

// получить данные по аксес токену
func TestGetDataByAccess(t *testing.T) {

	tests := []struct {
		description string
		expected    error
		token       string
	}{
		{
			description: "get data by access",
			expected:    nil,
			token:       createTokenAccess("8e2b4136-829c-11eb-8dcd-0242ac130003", "127.0.0.1"),
		},
	}

	for _, test := range tests {

		data, resp := getAccessToken(test.token)
		assert.Equalf(t, test.expected, resp, test.description)
		assert.Equalf(t, "8e2b4136-829c-11eb-8dcd-0242ac130003", data.GUID, test.description)

	}
}

// получить данные по рефреш токену
func TestGetDataByRefresh(t *testing.T) {

	tests := []struct {
		description string
		expected    error
		token       string
	}{
		{
			description: "get data by access",
			expected:    nil,
			token:       createTokenRefresh("8e2b4136-829c-11eb-8dcd-0242ac130003", createTokenAccess("8e2b4136-829c-11eb-8dcd-0242ac130003", "127.0.0.1")),
		},
	}

	for _, test := range tests {

		data, resp := getRefreshToken(test.token)
		assert.Equalf(t, test.expected, resp, test.description)
		assert.Equalf(t, "8e2b4136-829c-11eb-8dcd-0242ac130003", data.GUID, test.description)

	}
}

// дополнительная функция для тестов. создаёт просроченый токен.
func oldTokenCreate() string {
	var tokenObj = models.TokenAccessData{
		GUID:  "8e2b4136-829c-11eb-8dcd-0242ac130003",
		Ip:    "127.0.0.1",
		Email: GlobalEmail.Reciever,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(0 * time.Microsecond).Unix(),
		},
	}

	// method HS = HMAC + SHA 512
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, tokenObj)
	tokenString, err := token.SignedString(accessKey)
	if err != nil {
		log.Error("token dont signed")
		return ""
	}

	return tokenString
}
