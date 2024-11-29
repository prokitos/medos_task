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

func (instance ResponseBase) GetError() error {
	var temp = "Code:" + instance.Code + " ; " + " Description:" + instance.Description
	return errors.New(temp)
}
