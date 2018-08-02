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

func TestCreateANewDockerRegistryWithoutAuth(t *testing.T) {
	// given

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	registry := &model.DockerRegistry{Name: "test", Url: "localhost:5000", User: "tester", Password: "test"}
	body, _ := json.Marshal(registry)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/containers/docker/registries",
		s.URL), bytes.NewBuffer(body))
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)
	// shutdown server and db gracefully
	s.Close()
	tests.CleanPersistence(conf)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 401 {
		t.Fatalf("Expect 401 return code when trying to create a docker registry without auth. "+
			"Got %d", resp.StatusCode)
	}
}

func TestCreateANewDockerRegistry(t *testing.T) {
	// given

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	registry := &model.DockerRegistry{Name: "test", Url: "localhost:5000", User: "tester", Password: "test"}
	body, _ := json.Marshal(registry)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/containers/docker/registries",
		s.URL), bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)
	// shutdown server and db gracefully
	s.Close()
	tests.CleanPersistence(conf)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 201 {
		t.Fatalf("Expect 201 return code when trying to create a docker registry. "+
			"Got %d", resp.StatusCode)
	}
	var newRegistry model.DockerRegistry
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&newRegistry)
	if newRegistry.Name != registry.Name {
		t.Fatalf("expected Name %s from file to equals %s", newRegistry.Name, registry.Name)
	}
	if newRegistry.User != "" {
		t.Fatalf("expected User to have been removed but got %s", newRegistry.User)
	}
	if newRegistry.Password != "" {
		t.Fatalf("expected Password to have been removed but got %s", newRegistry.Password)
	}
}

func TestCheckPrivateRegistryNotAuthenticated(t *testing.T) {
	// given
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	registry := &model.DockerRegistry{Name: "test", Url: "localhost:5000", User: "tester", Password: "test"}
	body, _ := json.Marshal(registry)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/containers/docker/registries/test",
		s.URL), bytes.NewBuffer(body))
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have no error, but got %s", err.Error())
	}
	if resp.StatusCode != 401 {
		t.Fatalf("Expect 401 return code when testing a private docker registry without been previously authenticated"+
			"Got %d", resp.StatusCode)
	}
}

func TestCheckPrivateRegistrySuccessfully(t *testing.T) {
	// given
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	registry := &model.DockerRegistry{Name: "test", Url: "localhost:5000", User: "tester", Password: "test"}
	body, _ := json.Marshal(registry)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/containers/docker/registries/test",
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
		t.Fatalf("Expect 200 return code when testing successfully a private docker registry with user / password "+
			"Got %d", resp.StatusCode)
	}
}

func TestCheckRegistryBadDockerCredential(t *testing.T) {
	// given
	expectedErrorMsg := "Error response from daemon: login attempt to http://localhost:5000/v2/ failed with status: 401 Unauthorized"
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	registry := &model.DockerRegistry{Name: "test", Url: "localhost:5000", User: "tester", Password: "bad password"}
	body, _ := json.Marshal(registry)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/containers/docker/registries/test",
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
		t.Fatalf("Expect 403 return code when testing an private docker registry with wrong password "+
			"Got %d", resp.StatusCode)
	}
	var apiErr api.ApiError
	d := json.NewDecoder(resp.Body)
	d.Decode(&apiErr)
	if apiErr.Msg != expectedErrorMsg {
		t.Fatalf("Expect %s message when testing an unknown private docker registry with user / password "+
			"Got %s", expectedErrorMsg, apiErr.Msg)
	}
}

func TestCheckRegistryNoDockerCredential(t *testing.T) {
	// given
	expectedErrorMsg := "Error response from daemon: Get http://localhost:5000/v2/: no basic auth credentials"
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	registry := &model.DockerRegistry{Name: "test", Url: "localhost:5000", User: "", Password: ""}
	body, _ := json.Marshal(registry)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/containers/docker/registries/test",
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
		t.Fatalf("Expect 403 return code when testing an private docker registry with wrong password "+
			"Got %d", resp.StatusCode)
	}
	var apiErr api.ApiError
	d := json.NewDecoder(resp.Body)
	d.Decode(&apiErr)
	if apiErr.Msg != expectedErrorMsg {
		t.Fatalf("Expect %s message when testing an unknown private docker registry with user / password "+
			"Got %s", expectedErrorMsg, apiErr.Msg)
	}
}

func TestCheckUnknownRegistry(t *testing.T) {
	// given
	expectedErrorMsg := "Error response from daemon: Get https://hotelocal:5000/v2/: dial tcp: lookup hotelocal: no such host"
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// server setup
	a := api.InitRoute(conf)
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(a.Handler())
	defer s.Close()

	// request setup
	registry := &model.DockerRegistry{Name: "test", Url: "hotelocal:5000", User: "tester", Password: "test"}
	body, _ := json.Marshal(registry)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/containers/docker/registries/test",
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
		t.Fatalf("Expect 404 return code when testing an unknown private docker registry with user / password "+
			"Got %d", resp.StatusCode)
	}
	var apiErr api.ApiError
	d := json.NewDecoder(resp.Body)
	d.Decode(&apiErr)
	if apiErr.Msg != expectedErrorMsg {
		t.Fatalf("Expect %s message when testing an unknown private docker registry with user / password "+
			"Got %s", expectedErrorMsg, apiErr.Msg)
	}
}
