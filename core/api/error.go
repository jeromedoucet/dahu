package api

import "encoding/json"

type apiError struct {
	Msg string `json:"msg"`
}

func fromErrorToJson(err error) []byte {
	// todo test nil !
	apiErr := apiError{Msg: err.Error()}
	res, _ := json.Marshal(apiErr) // todo handle err
	return res
}
