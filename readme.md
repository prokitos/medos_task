Сборка приложения:

вариант 1 )  go build cmd/main.go

вариант 2 )  make

--------

Сервер по стандарту запущен локально на 8001 порту.

GET роут /getToken      принимает [GUID]

GET роут /refreshToken  принимает [access] и [refresh]

--------

Для запуска в докере нужно в [main.go] на 15 строке заменить значение переменной на [var CONFIG_TYPE string = "docker"]

после этого:

вариант 1 ) вручную билдить докер образ через композ 

вариант 2 ) поменять 7 строку в makefile на [currentDepoly: deployDocker] ;  дальше выполнить команду make


--------

Почта и пароль указываются в yaml файлах.

Конфигурация изменяется в [config/docker.yaml] или [config/local.yaml] 


