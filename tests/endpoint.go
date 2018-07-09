package tests

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func GetToken(secret string, exp time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": exp.Unix(),
	})
	res, _ := token.SignedString([]byte(secret))
	return res
}
