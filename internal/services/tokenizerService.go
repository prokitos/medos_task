package services

import (
	"errors"
	"mymod/internal/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

var accessKey = []byte("basic_key")
var refreshKey = []byte("super_secret_key")

// создание аксес токена.
func createTokenAccess(GUID string, ip string) string {

	// создаем токен, для теста укажим email при создании
	var tokenObj = models.TokenAccessData{
		GUID:  GUID,
		Ip:    ip,
		Email: "test@gmail.com",
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

// получить новый рефреш токен
func RenewToken(guid string, ip string) models.Tokens {

	var result models.Tokens

	// ВСЁ ВАЛИДНО! создание нового рефреш и аксес токена
	result.AccessToken = createTokenAccess(guid, ip)
	result.RefreshToken = createTokenRefresh(guid, result.AccessToken)

	// возвращаем токены обратно
	return result
}

func CheckTokens(refreshToken string, accessToken string) error {

	var errorr string

	// проверка рефреш токена
	token, err := validateRefreshToken(refreshToken)
	if err != nil {
		errorr = "refresh token unauthorized"
		if err == jwt.ErrSignatureInvalid {
			errorr = "refresh token sign unknown"
			return errors.New(errorr)
		}
		return errors.New(errorr)
	}

	if !token.Valid {
		errorr = "refresh token expired"
		return errors.New(errorr)
	}

	// проверяем что аксес токен внутри рефреш токена совпадает с нашим аксес токеном.
	refToken := token.Claims.(*models.TokenRefreshData)
	if refToken.AcceessToken != accessToken {
		errorr = "access token missmatch"
		return errors.New(errorr)
	}

	return nil
}

func GetGuid(access string) (string, error) {

	token, err := validateRefreshToken(access)
	if err != nil {
		return "", err
	}

	accToken := token.Claims.(*models.TokenAccessData)
	return accToken.GUID, nil
}

func CheckTokenIp(refresh string, ip string) error {

	token, err := validateRefreshToken(refresh)
	if err != nil {
		return err
	}

	refToken := token.Claims.(*models.TokenRefreshData)

	token, err = validateRefreshToken(refToken.AcceessToken)
	if err != nil {
		return err
	}

	accToken := token.Claims.(*models.TokenAccessData)
	if accToken.Ip != ip {
		SendEmail(accToken.Email)
	}

	return nil
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
