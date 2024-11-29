package services

import (
	"mymod/internal/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

var accessKey = []byte("basic_key")
var refreshKey = []byte("super_secret_key")

// создание аксес токена.
func createTokenAccess(GUID string) string {

	// создаем токен
	var tokenObj = models.TokenAccessData{
		GUID: GUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
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

// создание рефреш токена.
func createTokenRefresh(GUID string, accessToken string) string {

	// создаем токен
	var tokenObj = models.TokenRefreshData{
		GUID:         GUID,
		AcceessToken: accessToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}

	// method HS = HMAC + SHA 512
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, tokenObj)
	tokenString, err := token.SignedString(refreshKey)

	if err != nil {
		log.Error("token dont signed")
		return ""
	}

	return tokenString
}

// проверка валидности access токена
func validateAccessToken(bearerToken string) (*jwt.Token, error) {

	tokenString := strings.Split(bearerToken, " ")[1]
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenAccessData{}, func(token *jwt.Token) (interface{}, error) {
		return accessKey, nil
	})
	return token, err
}

// проверка валидности refresh токена
func validateRefreshToken(bearerToken string) (*jwt.Token, error) {

	//tokenString := strings.Split(bearerToken, " ")[1]

	tokenString := bearerToken
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenRefreshData{}, func(token *jwt.Token) (interface{}, error) {
		return refreshKey, nil
	})

	return token, err
}
