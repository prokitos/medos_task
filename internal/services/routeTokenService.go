package services

import (
	"mymod/internal/database"
	"mymod/internal/models"
	"os"
	"strconv"
)

// получает рефреш и аксес токены
func RouteGetToken(guid string) (models.Tokens, error) {

	// получает ip.
	ip, err := os.Hostname()
	if err != nil {
		return models.Tokens{}, err
	}

	// создаёт пару токенов.
	var access = createTokenAccess(guid, ip)
	var refresh = createTokenRefresh(guid, access)

	// так как в таблице хранятся только токены, то просто удаляется вся запись с токенами для этого guid.
	database.GlobalPostgres.DeleteData(models.Auth{GUID: guid})

	// создаёт записи в таблице.
	var newRecord models.Auth
	newRecord.GUID = guid
	newRecord.Refresh = refresh
	err = database.GlobalPostgres.CreateData(newRecord)
	if err != nil {
		return models.Tokens{}, err
	}

	// формируется результат.
	var result models.Tokens
	result.AccessToken = access
	result.RefreshToken = refresh
	return result, nil
}

// проверяет что рефреш и аксес токены валидны, выдаём новые токены
func RouteRefreshToken(access string, refresh string) (models.Tokens, error) {

	// проверяет валидность токенов. также проверяем что рефреш токен связан с аксес токеном.
	err := checkAccessToken(access)
	if err != nil {
		return models.Tokens{}, err
	}
	err = checkRefreshTokens(refresh, access)
	if err != nil {
		return models.Tokens{}, err
	}

	// проверяет что наш рефреш токен есть в базе.
	var curUser models.Auth
	curUser.Refresh = refresh
	err = database.GlobalPostgres.CheckExist(curUser)
	if err != nil {
		return models.Tokens{}, err
	}

	// проверка ip. если не совпадают то отправка сообщения на почту.
	ip, err := os.Hostname()
	if err != nil {
		return models.Tokens{}, err
	}
	err = checkTokenIp(refresh, ip)
	if err != nil {
		return models.Tokens{}, err
	}

	// создаёт новые токены
	guid, err := getGuid(access)
	if err != nil {
		return models.Tokens{}, err
	}
	newToken := renewToken(guid, ip)

	// обновляет токены в базе. сначала получает id записи по старому refreshToken, а потом по id меняет refreshToken на новый.
	var newData models.Auth
	strid, err := database.GlobalPostgres.GetId(models.Auth{Refresh: refresh})
	if err != nil {
		return models.Tokens{}, err
	}
	id, err := strconv.Atoi(strid)
	if err != nil {
		return models.Tokens{}, err
	}
	newData.UserId = id
	newData.Refresh = newToken.RefreshToken
	database.GlobalPostgres.UpdateData(newData)

	return newToken, nil
}
