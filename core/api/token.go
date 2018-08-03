package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func (a *Api) authFilter(w http.ResponseWriter, r *http.Request) bool {
	err := a.checkToken(r)
	if err == nil {
		return true
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
}

// check if the given request contains a valid JWT token
// This function will search the token in the authorization
// header. The supported authentication scheme is bearer.
// it means that the expected header is Authorization : Bearer <TOKEN>
func (a *Api) checkToken(r *http.Request) (err error) {
	authContent := r.Header.Get("Authorization")
	chunck := strings.Split(strings.TrimSpace(authContent), " ")
	if len(chunck) != 2 {
		err = errors.New("invalid authorization data. Must have the form Bearer <TOKEN>")
		return
	}
	token, parsingError := jwt.Parse(chunck[1], a.keyFunc)
	if parsingError != nil {
		err = parsingError
		return
	}
	if !token.Valid {
		err = errors.New("invalid token")
		return
	}
	return
}

func (a *Api) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(a.conf.ApiConf.Secret), nil
}
