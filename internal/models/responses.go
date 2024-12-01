package models

import "errors"

type ResponseBase struct {
	Description string `json:"description" example:"description"`
}

func (instance ResponseBase) BaseServerError() error {
	var temp ResponseBase
	temp.Description = "Internal Error"
	instance = temp
	return instance.GetError()
}

func (instance ResponseBase) BadCreate() error {
	var temp ResponseBase
	temp.Description = "Create error"
	instance = temp
	return instance.GetError()
}

func (instance ResponseBase) BadUpdate() error {
	var temp ResponseBase
	temp.Description = "Update error"
	instance = temp
	return instance.GetError()
}

func (instance ResponseBase) BadDelete() error {
	var temp ResponseBase
	temp.Description = "Delete error"
	instance = temp
	return instance.GetError()
}

func (instance ResponseBase) BadShow() error {
	var temp ResponseBase
	temp.Description = "Data not exist"
	instance = temp
	return instance.GetError()
}

func (instance ResponseBase) CustomTokenError(text string) error {
	var temp ResponseBase
	temp.Description = text
	instance = temp
	return instance.GetError()
}

func (instance ResponseBase) GetError() error {
	return errors.New(instance.Description)
}
