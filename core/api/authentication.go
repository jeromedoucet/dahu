package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jeromedoucet/dahu/core/model"
)

func (a *Api) handleAuthentication(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	l := model.Login{}
	d := json.NewDecoder(r.Body)
	d.Decode(&l) // todo handle this error
	u, err := a.repository.GetUser(l.Id, ctx)
	if err != nil {
		log.Printf("INFO >> handleAuthentication unknown user id %v for authentication", string(l.Id))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = u.ComparePassword([]byte(l.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token := createToken(a.conf.ApiConf.Secret, time.Now().Add(a.conf.ApiConf.TokenValidityDuration))
	res := model.Token{Value: token}
	body, _ := json.Marshal(res) // todo handle err
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", body)
}

func createToken(secret string, exp time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": exp.Unix(),
	})
	res, _ := token.SignedString([]byte(secret))
	return res
}
