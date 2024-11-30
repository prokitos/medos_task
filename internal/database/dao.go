package database

import (
	"mymod/internal/models"

	log "github.com/sirupsen/logrus"
)

func (currentlDB *PostgresDatabase) CreateData(data models.Auth) error {

	if result := currentlDB.Instance.Create(&data); result.Error != nil {
		log.Debug("create record error!")
		return models.ResponseBase{}.BadCreate()
	}

	log.Debug("dao complete")
	return models.ResponseBase{}.GoodCreate()
}

func (currentlDB *PostgresDatabase) CheckExist(data models.Auth) error {

	var finded []models.Auth

	results := currentlDB.Instance.Find(&finded, data)
	if results.Error != nil || results.RowsAffected == 0 {
		log.Debug("show record error!")
		return models.ResponseBase{}.BadShow()
	}

	log.Debug("dao complete")
	return nil
}

func (currentlDB *PostgresDatabase) UpdateData(data models.Auth) error {

	if result := currentlDB.Instance.Updates(&data); result.Error != nil {
		log.Debug("update record error!")
		return models.ResponseBase{}.BadUpdate()
	}

	log.Debug("dao complete")
	return models.ResponseBase{}.GoodUpdate()
}

func (currentlDB *PostgresDatabase) GetId(data models.Auth) string {

	return ""
}
