package models

import "github.com/golang-jwt/jwt/v4"

// модели для работы с токенами
type TokenAccessData struct {
	GUID  string
	Ip    string
	Email string
	jwt.StandardClaims
}

type TokenRefreshData struct {
	GUID         string
	AcceessToken string
	jwt.StandardClaims
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
