package services

import (
	"mymod/internal/database"
	"mymod/internal/models"
	"net"
	"strconv"
)

// получить айпи
func resolveHostIp() string {

	netInterfaceAddresses, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, netInterfaceAddress := range netInterfaceAddresses {
		networkIp, ok := netInterfaceAddress.(*net.IPNet)
		if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {
			ip := networkIp.IP.String()
			return ip
		}
	}
	return ""
}

// получает рефреш и аксес токены
func RouteGetToken(guid string) (models.Tokens, error) {

	// получает ip.
	ip := resolveHostIp()

	// создаёт пару токенов.
	var access = createTokenAccess(guid, ip)
	var refresh = createTokenRefresh(guid, access)

	// так как в таблице хранятся только токены, то просто удаляется вся запись с токенами для этого guid.
	database.GlobalPostgres.DeleteDataByGuid(models.Auth{GUID: guid})

	// создаёт записи в таблице.
	var newRecord models.Auth
	newRecord.GUID = guid
	newRecord.Refresh = refresh
	err := database.GlobalPostgres.CreateData(newRecord)
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
	ip := resolveHostIp()
	checkTokenIp(refresh, ip)

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
