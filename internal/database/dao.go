package database

import (
	"mymod/internal/models"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (currentlDB *PostgresDatabase) CreateData(data models.Auth) error {

	if result := currentlDB.Instance.Create(&data); result.Error != nil {
		log.Debug("create record error!")
		return models.ResponseBase{}.BadCreate()
	}

	log.Debug("create complete")
	return nil
}

func (currentlDB *PostgresDatabase) UpdateData(data models.Auth) error {

	if result := currentlDB.Instance.Updates(&data); result.Error != nil {
		log.Debug("update record error!")
		return models.ResponseBase{}.BadUpdate()
	}

	log.Debug("update complete")
	return nil
}

func (currentlDB *PostgresDatabase) DeleteData(data models.Auth) error {

	// result := currentlDB.Instance.Delete(&data, data.GUID)
	result := currentlDB.Instance.Model(&models.Auth{}).Where("guid= ?", data.GUID).Delete(&data)
	if result.RowsAffected == 0 || result.Error != nil {
		return models.ResponseBase{}.BadDelete()
	}

	log.Debug("delete complete")
	return nil

}

func (currentlDB *PostgresDatabase) CheckExist(data models.Auth) error {

	var finded models.Auth

	results := currentlDB.Instance.Find(&finded, data)
	if results.Error != nil || results.RowsAffected == 0 {
		log.Debug("show record error!")
		return models.ResponseBase{}.BadShow()
	}

	log.Debug("check complete")
	return nil
}

func (currentlDB *PostgresDatabase) GetId(data models.Auth) (string, error) {

	var finded models.Auth

	results := currentlDB.Instance.Find(&finded, data)
	if results.Error != nil || results.RowsAffected == 0 {
		log.Debug("show record error!")
		return "", models.ResponseBase{}.BadShow()
	}

	log.Debug("get complete")
	strid := strconv.Itoa(finded.UserId)
	return strid, nil
}
