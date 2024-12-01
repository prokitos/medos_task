package services

import (
	"mymod/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

var accessKey = []byte("basic_key")
var refreshKey = []byte("super_secret_key")

// создание аксес токена.
func createTokenAccess(GUID string, ip string) string {

	// создает токен, для теста email указан при создании.
	var tokenObj = models.TokenAccessData{
		GUID:  GUID,
		Ip:    ip,
		Email: GlobalEmail.Reciever,
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

	// создает токен
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

	//tokenString := strings.Split(bearerToken, " ")[1]
	tokenString := bearerToken
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

// получает данные по аксес токену
func getAccessToken(access string) (*models.TokenAccessData, error) {
	token, err := validateAccessToken(access)
	if err != nil {
		return nil, err
	}

	data := token.Claims.(*models.TokenAccessData)
	return data, nil
}

// получает данные по рефреш токену
func getRefreshToken(refresh string) (*models.TokenRefreshData, error) {
	token, err := validateRefreshToken(refresh)
	if err != nil {
		return nil, err
	}

	data := token.Claims.(*models.TokenRefreshData)
	return data, nil
}

// получает guid по аксес токену
func getGuid(access string) (string, error) {

	data, err := getAccessToken(access)
	if err != nil {
		return "", err
	}

	return data.GUID, nil
}

// проверка совпадает ли текущий айпи и айпи в рефреш токене
func checkTokenIp(refresh string, ip string) error {

	data, err := getRefreshToken(refresh)
	if err != nil {
		return err
	}

	newdata, err := getAccessToken(data.AcceessToken)
	if err != nil {
		return err
	}

	if newdata.Ip != ip {
		SendEmail(newdata.Email)
	}

	return nil
}

// получить новый рефреш токен
func renewToken(guid string, ip string) models.Tokens {

	var result models.Tokens

	// ТУТ УЖЕ ДОЛЖНО БЫТЬ ВСЁ ВАЛИДНО и поэтому идёт просто создание нового рефреш и аксес токена
	result.AccessToken = createTokenAccess(guid, ip)
	result.RefreshToken = createTokenRefresh(guid, result.AccessToken)

	// возвращаем токены обратно
	return result
}

// проверка рефреш токена
func checkRefreshTokens(refreshToken string, accessToken string) error {

	var errorr string
	token, err := validateRefreshToken(refreshToken)
	if err != nil {
		errorr = "refresh token unauthorized"
		if err == jwt.ErrSignatureInvalid {
			errorr = "refresh token sign unknown"
			return models.ResponseBase{}.CustomTokenError(errorr)
		}
		return models.ResponseBase{}.CustomTokenError(errorr)
	}

	if !token.Valid {
		errorr = "refresh token expired"
		return models.ResponseBase{}.CustomTokenError(errorr)
	}

	// проверяем что аксес токен внутри рефреш токена совпадает с нашим аксес токеном.
	refToken := token.Claims.(*models.TokenRefreshData)
	if refToken.AcceessToken != accessToken {
		errorr = "access token missmatch"
		return models.ResponseBase{}.CustomTokenError(errorr)
	}

	return nil
}

// проверка аксес токена
func checkAccessToken(accessToken string) error {

	var errorr string
	token, err := validateAccessToken(accessToken)
	if err != nil {
		errorr = "access token unauthorized"
		if err == jwt.ErrSignatureInvalid {
			errorr = "access token sign unknown"
			return models.ResponseBase{}.CustomTokenError(errorr)
		}
		return models.ResponseBase{}.CustomTokenError(errorr)
	}

	if !token.Valid {
		errorr = "access token expired"
		return models.ResponseBase{}.CustomTokenError(errorr)
	}

	return nil
}
