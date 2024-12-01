package database

import (
	"fmt"
	"mymod/internal/config"
	"mymod/internal/models"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// загрузка тестовых конфигов
func getTestConfig() config.PostgresConfig {

	var cfg config.PostgresConfig
	cfg.Host = "localhost"
	cfg.Name = "meados_test"
	cfg.Pass = "root"
	cfg.Port = "8092"
	cfg.User = "postgres"
	cfg.Reload = false

	return cfg
}

// проверка подключения к базе. создание тестовой базы данных.
func TestConnetion(t *testing.T) {

	var cfg config.MainConfig
	cfg.PostgresDB = getTestConfig()

	var temp PostgresDatabase
	temp.Run(cfg)

	sqlDB, err := temp.Instance.DB()
	if err != nil {
		t.Errorf(err.Error())
	}

	qq := sqlDB.Stats()
	if qq.Idle == 0 && qq.InUse == 0 {
		t.Errorf("result wrong at test, does not connect to server")
	}

	sqlDB.Close()
}

// создаём запись в тестовой таблице.
func TestCreateData(t *testing.T) {

	var cfg config.MainConfig
	cfg.PostgresDB = getTestConfig()

	var temp PostgresDatabase
	temp.openConnection(cfg)

	var data models.Auth
	data.GUID = "3f2504e0-4f89-11d3-9a0c-0305e82c3301"
	data.Refresh = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJHVUlEIjoiZ2VyLWhlcm1lbiIsIklwIjoiMjQzZTU0NzExY2ZmIiwiRW1haWwiOiJ3d3d3QGdtYWlsLmNvbSIsImV4cCI6MTczMzA1NTExOH0.RdbPRUKIj0G8XTirifKLdYIb1CW89zpYswC3jLAetDuc3pJ_p6D7GrXUoKJngmN0neHOoJAceov4rbYw2LAI6w"
	err := temp.CreateData(data)

	if err != nil {
		t.Errorf("error, data don't create")
	}

}

// проверяем есть ли запись с нужным guid
func TestExistData(t *testing.T) {

	var cfg config.MainConfig
	cfg.PostgresDB = getTestConfig()

	var temp PostgresDatabase
	temp.openConnection(cfg)

	var data models.Auth
	data.GUID = "3f2504e0-4f89-11d3-9a0c-0305e82c3301"
	err := temp.CheckExist(data)

	if err != nil {
		t.Errorf("error, data don't exist")
	}

}

// обновляем токен у записи
func TestUpdateData(t *testing.T) {

	var cfg config.MainConfig
	cfg.PostgresDB = getTestConfig()

	var temp PostgresDatabase
	temp.openConnection(cfg)

	var data models.Auth
	data.UserId = 1
	data.GUID = "3f2504e0-4f89-11d3-9a0c-0305e82c3301"
	data.Refresh = "XTirifKLdYIb1CW89zpYswC3jLAetDuc3pJ_p6D7GrXUoKJngmN0neHOoJAceov4rbYw2LAI6w"
	err := temp.UpdateData(data)

	if err != nil {
		t.Errorf("error, data don't update")
	}

}

// получаем айди обновлённой записи по токену.
func TestGetId(t *testing.T) {

	var cfg config.MainConfig
	cfg.PostgresDB = getTestConfig()

	var temp PostgresDatabase
	temp.openConnection(cfg)

	var data models.Auth
	data.Refresh = "XTirifKLdYIb1CW89zpYswC3jLAetDuc3pJ_p6D7GrXUoKJngmN0neHOoJAceov4rbYw2LAI6w"
	newdata, err := temp.GetId(data)

	if err != nil && newdata != "1" {
		t.Errorf("error, id don't get correct")
	}

}

// удаляем запись по гуиду.
func TestDeleteData(t *testing.T) {

	var cfg config.MainConfig
	cfg.PostgresDB = getTestConfig()

	var temp PostgresDatabase
	temp.openConnection(cfg)

	var data models.Auth
	data.GUID = "3f2504e0-4f89-11d3-9a0c-0305e82c3301"
	err := temp.DeleteDataByGuid(data)

	if err != nil {
		t.Errorf("error, data don't delete")
	}

}

// проверяем что записей у этого гуида не осталось.
func TestExistDataSecond(t *testing.T) {

	var cfg config.MainConfig
	cfg.PostgresDB = getTestConfig()

	var temp PostgresDatabase
	temp.openConnection(cfg)

	var data models.Auth
	data.GUID = "3f2504e0-4f89-11d3-9a0c-0305e82c3301"
	err := temp.CheckExist(data)

	if err == nil {
		t.Errorf("error, data exist after delete")
	}

}

// удаляем тестовую таблицу после тестов.
func TestEnd(t *testing.T) {

	var cfg config.MainConfig
	cfg.PostgresDB = getTestConfig()

	connectStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PostgresDB.User, cfg.PostgresDB.Pass, cfg.PostgresDB.Host, cfg.PostgresDB.Port, cfg.PostgresDB.Name)
	db, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{})
	if err != nil {
		t.Errorf("database dont open")
	}

	// закрытие бд
	sql, _ := db.DB()
	defer func() {
		_ = sql.Close()
	}()

	stmt := fmt.Sprintf("DROP TABLE %s;", "auths")
	if rs := db.Exec(stmt); rs.Error != nil {
		t.Errorf("table dont delete")
	}

}
