package services

import (
	"mymod/internal/database"
	"mymod/internal/models"
	"os"
)

func RouteGetToken(guid string) (models.Tokens, error) {

	ip, err := os.Hostname()
	if err != nil {
		return models.Tokens{}, err
	}

	var access = createTokenAccess(guid, ip)
	var refresh = createTokenRefresh(guid, access)

	var newRecord models.Auth
	newRecord.GUID = guid
	newRecord.Refresh = refresh
	err = database.GlobalPostgres.CreateData(newRecord)
	if err != nil {
		return models.Tokens{}, err
	}

	var result models.Tokens
	result.AccessToken = access
	result.RefreshToken = refresh
	return result, nil
}

func RouteRefreshToken(access string, refresh string) (models.Tokens, error) {

	// проверяем что наш рефреш токен есть в базе.
	var curUser models.Auth
	curUser.Refresh = refresh
	err := database.GlobalPostgres.CheckExist(curUser)
	if err != nil {
		return models.Tokens{}, err
	}

	// проверка ip
	ip, err := os.Hostname()
	if err != nil {
		return models.Tokens{}, err
	}
	err = CheckTokenIp(refresh, ip)
	if err != nil {
		return models.Tokens{}, err
	}

	// проверяем что аксес токен совпадает с рефреш токеном
	err = CheckTokens(refresh, access)
	if err != nil {
		return models.Tokens{}, err
	}

	// создаём новые токены
	guid, err := GetGuid(access)
	if err != nil {
		return models.Tokens{}, err
	}
	newToken := RenewToken(guid, ip)

	// обновляем токены в базе
	var newData models.Auth
	newData.UserId = database.GlobalPostgres.GetId(models.Auth{Refresh: refresh})
	newData.Refresh = newToken.RefreshToken
	database.GlobalPostgres.UpdateData(newData)

	return newToken, nil

}
