package database

import (
	"fmt"
	"mymod/internal/config"
	"mymod/internal/models"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GlobalPostgres *PostgresDatabase

type PostgresDatabase struct {
	Instance *gorm.DB
}

func (currentlDB *PostgresDatabase) Run(config config.MainConfig) {
	currentlDB.openConnection(config)
	currentlDB.startMigration()
	currentlDB.globalSet()
}

func (currentlDB *PostgresDatabase) startMigration() {
	currentlDB.Instance.AutoMigrate(models.Auth{})
	log.Debug("migration complete")
}

func (currentlDB *PostgresDatabase) openConnection(config config.MainConfig) {

	err := currentlDB.checkDatabaseCreated(config)
	if err != nil {
		panic("not connection to db")
	}

	connectStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.PostgresDB.User, config.PostgresDB.Pass, config.PostgresDB.Host, config.PostgresDB.Port, config.PostgresDB.Name)

	db, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{})
	if err != nil {
		panic("not connection to db")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("not connection to db")
	}

	sqlDB.SetMaxIdleConns(4)
	sqlDB.SetMaxOpenConns(8)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	currentlDB.Instance = db
}

// проверка если есть база данных. если нет, то создать.
func (currentlDB *PostgresDatabase) checkDatabaseCreated(config config.MainConfig) error {

	// открытие соеднение с базой по стандарту
	connectStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.PostgresDB.User, config.PostgresDB.Pass, config.PostgresDB.Host, config.PostgresDB.Port, "postgres")
	db, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{})
	if err != nil {
		log.Error("database don't open")
		return models.ResponseBase{}.BaseServerError()
	}

	// закрытие бд
	sql, _ := db.DB()
	defer func() {
		_ = sql.Close()
	}()

	// проверка если есть нужная нам база данных
	stmt := fmt.Sprintf("SELECT * FROM pg_database WHERE datname = '%s';", config.PostgresDB.Name)
	rs := db.Raw(stmt)
	if rs.Error != nil {
		log.Error("error, dont read bd")
		return models.ResponseBase{}.BaseServerError()
	}

	// если нет, то создать
	var rec = make(map[string]interface{})
	if rs.Find(rec); len(rec) == 0 {
		stmt := fmt.Sprintf("CREATE DATABASE %s;", config.PostgresDB.Name)
		if rs := db.Exec(stmt); rs.Error != nil {
			log.Error("error, dont create a database")
			return models.ResponseBase{}.BaseServerError()
		}
	}

	return nil
}

func (currentlDB *PostgresDatabase) globalSet() {
	GlobalPostgres = currentlDB
}
