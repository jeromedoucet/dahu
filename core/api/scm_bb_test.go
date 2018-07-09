package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/api"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/tests"
)

// todo verify ssh without password

func TestCheckWhenNotAuthenticated(t *testing.T) {
	t.Skip()
	// id = tester
	// email = tester@tester.org
	// password = test
	// repo = http://localhost:10080/tester/test-repo.git
}

func TestCheckPrivateRepoConfigurationHttpMissingParams(t *testing.T) {
	t.Skip()
	// id = tester
	// email = tester@tester.org
	// password = test
	// repo = http://localhost:10080/tester/test-repo.git
}

func TestCheckPrivateRepoConfigurationHttpBadCredentials(t *testing.T) {
	// given
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer a.Close()
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.HttpAuthConfig{Url: "http://localhost:10080/tester/test-repo.git", User: "tester", Password: "wrong-password"}
	gitConfig := model.GitConfig{HttpAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repositorie",
		s.URL), bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have no error, but got %s", err.Error())
	}
	if resp.StatusCode != 403 {
		t.Fatalf("Expect 403 return code when testing a private git repository with bad http credentials "+
			"Got %d", resp.StatusCode)
	}
}

func TestCheckPrivateRepoConfigurationHttpUnknowUrl(t *testing.T) {
	// given
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer a.Close()
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.HttpAuthConfig{Url: "http://localhost:10080/tester/unknown.git", User: "tester", Password: "test"}
	gitConfig := model.GitConfig{HttpAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repositorie",
		s.URL), bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have no error, but got %s", err.Error())
	}
	if resp.StatusCode != 404 {
		t.Fatalf("Expect 404 return code when testing an unknown git repository"+
			"Got %d", resp.StatusCode)
	}
}

func TestCheckPrivateRepoConfigurationHttpSuccessfully(t *testing.T) {
	// given
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer a.Close()
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.HttpAuthConfig{Url: "http://localhost:10080/tester/test-repo.git", User: "tester", Password: "test"}
	gitConfig := model.GitConfig{HttpAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repositorie",
		s.URL), bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have no error, but got %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expect 200 return code when testing a private git repository with http credentials "+
			"Got %d", resp.StatusCode)
	}
}

func TestCheckPublicRepoConfigurationHttpSuccessfully(t *testing.T) {
	// given
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer a.Close()
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.HttpAuthConfig{Url: "http://localhost:10080/tester/test-repo-pub.git"}
	gitConfig := model.GitConfig{HttpAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repositorie",
		s.URL), bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have no error, but got %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expect 200 return code when testing a public git repository without credentials "+
			"Got %d", resp.StatusCode)
	}
}
