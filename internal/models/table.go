package models

type Auth struct {
	UserId  string `json:"user_id,omitempty" example:"12"`
	GUID    string `json:"guid" example:""`
	Refresh string `json:"refresh" example:""`
}
