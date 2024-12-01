package models

type Auth struct {
	UserId  int    `json:"user_id" example:"12"  gorm:"unique;primaryKey;autoIncrement"`
	GUID    string `json:"guid" example:""`
	Refresh string `json:"refresh" example:""`
}
