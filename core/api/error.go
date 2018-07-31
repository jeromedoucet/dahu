package api

import "encoding/json"

type ApiError struct {
	Msg string `json:"msg"`
}

func fromErrorToJson(err error) []byte {
	// todo test nil !
	apiErr := ApiError{Msg: err.Error()}
	res, _ := json.Marshal(apiErr) // todo handle err
	return res
}
