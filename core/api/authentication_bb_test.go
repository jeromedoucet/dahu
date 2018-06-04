package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/api"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/tests"
)

// testing succefull authentication
func TestAuthenticationShouldReturn200AndAToken(t *testing.T) {
	// given

	// setup the conf
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// insert existing user in the db
	password := "test_test_test_test"
	u := model.User{Login: "test"}
	u.SetPassword([]byte(password))
	tests.InsertObject(conf, []byte("users"), []byte(u.Login), u)

	// start the server and setup the request
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	l := model.Login{Id: u.Login, Password: password}
	body, _ := json.Marshal(l)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/login",
		s.URL), bytes.NewBuffer(body))
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// close, remove the db and shutdown the server
	s.Close()
	tests.CleanPersistence(conf)
	// then

	// check the response code and error
	if err != nil {
		t.Errorf("expect to get no error when trying to authenticate with correct credentials, but got %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("expect status code 200 when trying to authenticate with correct credentials, but got %d", resp.StatusCode)
	}

	// check the token itself
	tok := model.Token{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&tok)

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.ApiConf.Secret), nil
	}
	token, parsingError := jwt.Parse(tok.Value, keyFunc)
	if parsingError != nil {
		t.Errorf("expect to have no error parsing the token, but got : %s", parsingError.Error())
	}
	if !token.Valid {
		t.Errorf("expect the token to be valid, but this %+v is not", token)
	}
}

// testing authentication when no user found
func TestAuthenticationShouldReturn404AndNoTokenWhenNoUserFound(t *testing.T) {
	// given

	// setup the conf
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// insert existing user in the db
	password := "test_test_test_test"
	u := model.User{Login: "test"}

	// start the server and setup the request
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	l := model.Login{Id: u.Login, Password: password}
	body, _ := json.Marshal(l)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/login",
		s.URL), bytes.NewBuffer(body))
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// close, remove the db and shutdown the server
	s.Close()
	tests.CleanPersistence(conf)
	// then

	// check the response code and error
	if err != nil {
		t.Errorf("expect to get no error when trying to authenticate without user found, but got %s", err.Error())
	}
	if resp.StatusCode != 404 {
		t.Errorf("expect status code 404 when trying to authenticate without user found , but got %d", resp.StatusCode)
	}
	b := bytes.Buffer{}
	b.ReadFrom(resp.Body)
	if b.String() != "" {
		t.Errorf("expect to have no answer when no user found but got %s", b.String())
	}
}

// testing succefull authentication
func TestAuthenticationShouldReturn401AndNoTokenWhenBadPassword(t *testing.T) {
	// given

	// setup the conf
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// insert existing user in the db
	password := "test_test_test_test"
	u := model.User{Login: "test"}
	u.SetPassword([]byte(password))
	tests.InsertObject(conf, []byte("users"), []byte(u.Login), u)

	// start the server and setup the request
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	l := model.Login{Id: u.Login, Password: "totototototototototototo"}
	body, _ := json.Marshal(l)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/login",
		s.URL), bytes.NewBuffer(body))
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// close, remove the db and shutdown the server
	s.Close()
	tests.CleanPersistence(conf)
	// then

	// check the response code and error
	if err != nil {
		t.Errorf("expect to get no error when trying to authenticate with bad credential, but got %s", err.Error())
	}
	if resp.StatusCode != 401 {
		t.Errorf("expect status code 404 when trying to authenticate with bad credential, but got %d", resp.StatusCode)
	}
	b := bytes.Buffer{}
	b.ReadFrom(resp.Body)
	if b.String() != "" {
		t.Errorf("expect to have no answer when bad credential but got %s", b.String())
	}
}
