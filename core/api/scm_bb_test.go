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
	"github.com/jeromedoucet/dahu/tests/ssh_keys"
)

func TestCheckWhenNotAuthenticated(t *testing.T) {
	// given
	conf = configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/test-repo.git", gitRepoIp), Key: ssh_keys.PrivateUnprotected}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repository",
		s.URL), bytes.NewBuffer(body))
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have no error, but got %s", err.Error())
	}
	if resp.StatusCode != 401 {
		t.Fatalf("Expect 401 return code when testing a git repository without been authenticated "+
			"Got %d", resp.StatusCode)
	}
}

func TestCheckPrivateRepoConfigurationSshWithMissingKey(t *testing.T) {
	// given
	conf = configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/test-repo.git", gitRepoIp)}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repository",
		s.URL), bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have no error, but got %s", err.Error())
	}
	if resp.StatusCode != 400 {
		t.Fatalf("Expect 400 return code when testing a private git repository without ssh private key "+
			"Got %d", resp.StatusCode)
	}
}

func TestCheckPrivateRepoConfigurationSshWithUnknownRepository(t *testing.T) {
	// given
	conf = configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/test-toto-repo.git", gitRepoIp), Key: ssh_keys.PrivateProtected, KeyPassword: "tester"}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repository",
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

func TestCheckPrivateRepoConfigurationSshWithBadCredentials(t *testing.T) {
	// given
	conf = configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/test-repo.git", gitRepoIp), Key: ssh_keys.PrivateBad}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repository",
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
		t.Fatalf("Expect 403 return code when testing a private git repository with unregistered ssh private key "+
			"Got %d", resp.StatusCode)
	}
}

func TestCheckPrivateRepoConfigurationSshWithPasswordUnSuccessfully(t *testing.T) {
	// given
	conf = configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/test-repo.git", gitRepoIp), Key: ssh_keys.PrivateProtected, KeyPassword: "wrong-key-password"}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repository",
		s.URL), bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have no error, but got %s", err.Error())
	}
	if resp.StatusCode != 400 {
		t.Fatalf("Expect 400 return code when testing a private git repository with protected ssh private key but wrong passsword "+
			"Got %d", resp.StatusCode)
	}
}

func TestCheckPrivateRepoConfigurationSshWithPasswordSuccessfully(t *testing.T) {
	// given
	conf = configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/test-repo.git", gitRepoIp), Key: ssh_keys.PrivateProtected, KeyPassword: "tester"}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repository",
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
		t.Fatalf("Expect 200 return code when testing a private git repository with protected ssh private key "+
			"Got %d", resp.StatusCode)
	}
}

func TestCheckPrivateRepoConfigurationSshWithoutPasswordSuccessfully(t *testing.T) {
	// given
	conf = configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/test-repo.git", gitRepoIp), Key: ssh_keys.PrivateUnprotected}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repository",
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
		t.Fatalf("Expect 200 return code when testing a private git repository with ssh private key "+
			"Got %d", resp.StatusCode)
	}
}

func TestCheckPrivateRepoConfigurationHttpBadCredentials(t *testing.T) {
	// given
	conf = configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.HttpAuthConfig{Url: fmt.Sprintf("http://%s:3000/tester/test-repo.git", gitRepoIp), User: "tester", Password: "wrong-password"}
	gitConfig := model.GitConfig{HttpAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repository",
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
	conf = configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.HttpAuthConfig{Url: fmt.Sprintf("http://%s:3000/tester/unknown.git", gitRepoIp), User: "tester", Password: "test"}
	gitConfig := model.GitConfig{HttpAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repository",
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
	conf = configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.HttpAuthConfig{Url: fmt.Sprintf("http://%s:3000/tester/test-repo.git", gitRepoIp), User: "tester", Password: "test"}
	gitConfig := model.GitConfig{HttpAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repository",
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
	conf = configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	authConfig := model.HttpAuthConfig{Url: fmt.Sprintf("http://%s:3000/tester/test-repo-pub.git", gitRepoIp)}
	gitConfig := model.GitConfig{HttpAuth: &authConfig}
	body, _ := json.Marshal(gitConfig)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/scm/git/repository",
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
