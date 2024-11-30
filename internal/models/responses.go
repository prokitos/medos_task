package models

import "errors"

type ResponseBase struct {
	Description string `json:"description"       		 example:"description"`
	Code        string `json:"code"               		 example:"status"`
}

func (instance ResponseBase) BaseServerError() error {
	var temp ResponseBase
	temp.Code = "400"
	temp.Description = "Internal Error"
	return instance.GetError()
}

func (instance ResponseBase) GoodCreate() error {
	var temp ResponseBase
	temp.Code = "200"
	temp.Description = "Create successful"
	return instance.GetError()
}

func (instance ResponseBase) BadCreate() error {
	var temp ResponseBase
	temp.Code = "400"
	temp.Description = "Create error"
	return instance.GetError()
}

func (instance ResponseBase) GoodUpdate() error {
	var temp ResponseBase
	temp.Code = "200"
	temp.Description = "Update successful"
	return instance.GetError()
}

func (instance ResponseBase) BadUpdate() error {
	var temp ResponseBase
	temp.Code = "400"
	temp.Description = "Update error"
	return instance.GetError()
}

func (instance ResponseBase) BadShow() error {
	var temp ResponseBase
	temp.Code = "400"
	temp.Description = "Data not exist"
	return instance.GetError()
}

func (instance ResponseBase) GetError() error {
	var temp = "Code:" + instance.Code + " ; " + " Description:" + instance.Description
	return errors.New(temp)
}
