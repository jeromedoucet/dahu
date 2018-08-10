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

// test correct rejection of non - sens requests on
// registries endpoint (a POST operation has not sens on theses kinds of endpoints)
func TestUnsuportedOperationOnRegistry(t *testing.T) {
	// given
	registry := &model.DockerRegistry{Name: "test", Url: "localhost:5000", User: "tester", Password: "test"}
	registry.GenerateId()
	registry.NewLastModificationTime()

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	tests.InsertObject(conf, []byte("dockerRegistries"), []byte(registry.Id), registry)
	defer tests.CleanPersistence(conf)

	// update changes
	registry.Name = "one-test"

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	defer s.Close()

	// request setup
	body, _ := json.Marshal(registry)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/containers/docker/registries/%s",
		s.URL, registry.Id), bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expect 404 return code when requesting an unsupported operation on registry endpoint. "+
			"Got %d", resp.StatusCode)
	}
}

// test correct rejection of non - sens requests on
// registries endpoint (delete has no sens here)
func TestUnsuportedOperationOnRegistries(t *testing.T) {
	// given

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	defer tests.CleanPersistence(conf)

	// update changes

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	defer s.Close()

	// request setup
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/containers/docker/registries",
		s.URL), nil)
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expect 404 return code when requesting an unsupported operation on registries endpoint. "+
			"Got %d", resp.StatusCode)
	}
}

// test updating a docker registry
// without auth token
func TestUpdateDockerRegistryNotAuthenticated(t *testing.T) {
	// given
	registry := new(model.DockerRegistryUpdate)
	registry.Name = "test"
	registry.Url = "localhost:5000"
	registry.User = "tester"
	registry.Password = "test"
	registry.GenerateId()
	registry.NewLastModificationTime()

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	defer tests.CleanPersistence(conf)

	// update changes

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	defer s.Close()

	// request setup
	body, _ := json.Marshal(registry)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/containers/docker/registries/%s",
		s.URL, registry.Id), bytes.NewBuffer(body))
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expect 401 return code when trying to update a docker registry without authentication. "+
			"Got %d", resp.StatusCode)
	}
}

// test updating a docker registry that
// does not exist
func TestUpdateUnknownDockerRegistry(t *testing.T) {
	// given
	registry := new(model.DockerRegistryUpdate)
	registry.Name = "test"
	registry.Url = "localhost:5000"
	registry.User = "tester"
	registry.Password = "test"
	registry.GenerateId()
	registry.NewLastModificationTime()

	expectedErrorMsg := fmt.Sprintf("No docker registry with id %s found", registry.Id)

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	defer tests.CleanPersistence(conf)

	// update changes

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	defer s.Close()

	// request setup
	body, _ := json.Marshal(registry)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/containers/docker/registries/%s",
		s.URL, registry.Id), bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expect 404 return code when trying to update an inexisting docker registry. "+
			"Got %d", resp.StatusCode)
	}
	var apiErr api.ApiError
	d := json.NewDecoder(resp.Body)
	d.Decode(&apiErr)
	if apiErr.Msg != expectedErrorMsg {
		t.Fatalf("Expect %s message when trying to update an unexisting docker registry. "+
			"Got %s", expectedErrorMsg, apiErr.Msg)
	}
}

// this case test a conflict when updating one docker registry.
// an optimistic lock mechanism based on the LastModificationTime
// field is used to detect read conflict. When such case happened,
// the data from db are return with a conflict status.
func TestUpdateDockerRegistryConflict(t *testing.T) {
	// given
	// NewLastModificationTime()
	registry := new(model.DockerRegistryUpdate)
	registry.Name = "test"
	registry.Url = "localhost:5000"
	registry.User = "tester"
	registry.Password = "test"
	registry.GenerateId()
	registry.NewLastModificationTime()

	referenceModificationTime := registry.LastModificationTime

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	tests.InsertObject(conf, []byte("dockerRegistries"), []byte(registry.Id), registry.DockerRegistry)
	defer tests.CleanPersistence(conf)

	// update changes
	registry.Name = "one-test"
	registry.NewLastModificationTime() // force the conflict here

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	defer s.Close()

	// request setup
	body, _ := json.Marshal(registry)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/containers/docker/registries/%s",
		s.URL, registry.Id), bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("Expect 409 return code when trying to update a docker registry with conflict. "+
			"Got %d", resp.StatusCode)
	}
	var updatedRegistry model.DockerRegistry
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&updatedRegistry)
	if updatedRegistry.Name != "test" {
		t.Fatalf("expected Name %s from file to equals %s", updatedRegistry.Name, "test")
	}
	if updatedRegistry.User != "tester" {
		t.Fatalf("expected User to have been perserved but got %s", updatedRegistry.User)
	}
	if updatedRegistry.Password != "" {
		t.Fatalf("expected Password to have been removed but got %s", updatedRegistry.Password)
	}
	if updatedRegistry.LastModificationTime != referenceModificationTime {
		t.Fatal("expected LastModificationTime to remain unchanged but it has change")
	}
}

// nominal test case for updating docker registry
func TestUpdateDockerRegistry(t *testing.T) {
	// given
	registry := new(model.DockerRegistryUpdate)
	registry.Name = "test"
	registry.Url = "localhost:5000"
	registry.User = "tester"
	registry.Password = "test"
	registry.GenerateId()
	registry.NewLastModificationTime()

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	tests.InsertObject(conf, []byte("dockerRegistries"), []byte(registry.Id), registry.DockerRegistry)
	defer tests.CleanPersistence(conf)

	// update changes
	registry.Name = "one-test"
	registry.ChangedFields = []string{"name"}

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	defer s.Close()

	// request setup
	body, _ := json.Marshal(registry)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/containers/docker/registries/%s",
		s.URL, registry.Id), bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expect 200 return code when trying to update a docker registry. "+
			"Got %d", resp.StatusCode)
	}
	var updatedRegistry model.DockerRegistry
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&updatedRegistry)
	if updatedRegistry.Name != registry.Name {
		t.Fatalf("expected Name %s from file to equals %s", updatedRegistry.Name, registry.Name)
	}
	if updatedRegistry.User != "tester" {
		t.Fatalf("expected User to have been perserved but got %s", updatedRegistry.User)
	}
	if updatedRegistry.Password != "" {
		t.Fatalf("expected Password to have been removed but got %s", updatedRegistry.Password)
	}
	if updatedRegistry.LastModificationTime <= registry.LastModificationTime {
		t.Fatal("expected LastModificationTime to have been update but it is not the case")
	}
}

// test case for docker registry list
// api endpoint when not authenticated
func TestListDockerRegistryNotAuthenticated(t *testing.T) {
	// given

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	defer tests.CleanPersistence(conf)

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	defer s.Close()

	// request setup
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/containers/docker/registries",
		s.URL), nil)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expect 401 return code when trying to list docker registries without authentication. "+
			"Got %d", resp.StatusCode)
	}
}

// Nominal test case for docker registry list
// api endpoint
func TestListDockerRegistry(t *testing.T) {
	// given
	registry1 := &model.DockerRegistry{Name: "test1", Url: "localhost:5001", User: "tester1", Password: "test1"}
	registry1.GenerateId()
	registry2 := &model.DockerRegistry{Name: "test2", Url: "localhost:5002", User: "tester2", Password: "test2"}
	registry2.GenerateId()

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	tests.InsertObject(conf, []byte("dockerRegistries"), []byte(registry1.Id), registry1)
	tests.InsertObject(conf, []byte("dockerRegistries"), []byte(registry2.Id), registry2)
	defer tests.CleanPersistence(conf)

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	defer s.Close()

	// request setup
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/containers/docker/registries",
		s.URL), nil)
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expect 200 return code when trying to list docker registries. "+
			"Got %d", resp.StatusCode)
	}

	var registryList []model.DockerRegistry
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&registryList)
	if len(registryList) != 2 {
		t.Fatalf("expect list docker registries to return 2 elements but got %d", len(registryList))
	}
	if registryList[0].Name != registry1.Name {
		t.Fatalf("expected Name %s from file to equals %s", registryList[0].Name, registry1.Name)
	}
	if registryList[0].User != "tester1" {
		t.Fatalf("expected User to have been perserved but got %s", registryList[0].User)
	}
	if registryList[0].Password != "" {
		t.Fatalf("expected Password to have been removed but got %s", registryList[0].Password)
	}
	if registryList[1].Name != registry2.Name {
		t.Fatalf("expected Name %s from file to equals %s", registryList[1].Name, registry2.Name)
	}
	if registryList[1].User != "tester2" {
		t.Fatalf("expected User to have been perserved but got %s", registryList[1].User)
	}
	if registryList[1].Password != "" {
		t.Fatalf("expected Password to have been removed but got %s", registryList[1].Password)
	}
}

// Test case for docker registry deletion
// api endpoint when not authenticated
func TestDeleteDockerRegistryNotAuthenticated(t *testing.T) {
	// given

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	defer tests.CleanPersistence(conf)

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/containers/docker/registries/%s",
		s.URL, "1"), nil)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)
	// shutdown server and db gracefully
	s.Close()

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expect 404 return code when trying to delete an docker registry when not authenticated. "+
			"Got %d", resp.StatusCode)
	}
}

// Test case for docker registry deletion
// api endpoint when resource doesn't exist
func TestDeleteUnknownDockerRegistry(t *testing.T) {
	// given
	expectedErrorMsg := "No docker registry with id 1 found"

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	defer tests.CleanPersistence(conf)

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	defer s.Close()

	// request setup
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/containers/docker/registries/%s",
		s.URL, "1"), nil)
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)
	// shutdown server and db gracefully

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expect 404 return code when trying to delete an unexisting docker registry. "+
			"Got %d", resp.StatusCode)
	}
	var apiErr api.ApiError
	d := json.NewDecoder(resp.Body)
	d.Decode(&apiErr)
	if apiErr.Msg != expectedErrorMsg {
		t.Fatalf("Expect %s message when trying to delete an unexisting docker registry. "+
			"Got %s", expectedErrorMsg, apiErr.Msg)
	}
}

// Nominal test case for docker registry deletion
// api endpoint
func TestDeleteDockerRegistry(t *testing.T) {
	// given
	registry := &model.DockerRegistry{Name: "test", Url: "localhost:5000", User: "tester", Password: "test"}
	registry.GenerateId()

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	tests.InsertObject(conf, []byte("dockerRegistries"), []byte(registry.Id), registry)
	defer tests.DeletePersistence(conf)

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/containers/docker/registries/%s",
		s.URL, registry.Id), nil)
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)
	// shutdown server and db gracefully
	s.Close()
	tests.ClosePersistence(conf)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expect 200 return code when trying to delete a docker registry. "+
			"Got %d", resp.StatusCode)
	}

	exist := tests.ObjectExist(conf, []byte("dockerRegistries"), []byte(registry.Id))

	if exist {
		t.Fatalf("Expect the docker registry to have been deleted, but it is still present")
	}
}

// nominal test case for getting one docker registry
func TestGetDockerRegistry(t *testing.T) {
	// given
	registry := &model.DockerRegistry{Name: "test", Url: "localhost:5000", User: "tester", Password: "test"}
	registry.GenerateId()

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	tests.InsertObject(conf, []byte("dockerRegistries"), []byte(registry.Id), registry)
	defer tests.CleanPersistence(conf)

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/containers/docker/registries/%s",
		s.URL, registry.Id), nil)
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)
	// shutdown server and db gracefully
	s.Close()

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expect 200 return code when trying to get a docker registry. "+
			"Got %d", resp.StatusCode)
	}
	var newRegistry model.DockerRegistry
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&newRegistry)
	if newRegistry.Name != registry.Name {
		t.Fatalf("expected Name %s from file to equals %s", newRegistry.Name, registry.Name)
	}
	if newRegistry.User != "tester" {
		t.Fatalf("expected User to have been perserved but got %s", newRegistry.User)
	}
	if newRegistry.Password != "" {
		t.Fatalf("expected Password to have been removed but got %s", newRegistry.Password)
	}
}

// test getting one docker registry without authentication
func TestGetDockerRegistryNotAuthenticated(t *testing.T) {
	// given

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	defer tests.CleanPersistence(conf)

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/containers/docker/registries/1",
		s.URL), nil)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)
	// shutdown server and db gracefully
	s.Close()

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expect 401 return code when trying to get a docker registry without auth. "+
			"Got %d", resp.StatusCode)
	}
}

// test getting an unexisting docker registry
func TestGetUnknownDockerRegistry(t *testing.T) {
	// given
	expectedErrorMsg := "No docker registry with id 1 found"

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	defer tests.CleanPersistence(conf)

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/containers/docker/registries/1",
		s.URL), nil)
	req.Header.Add("Authorization", "Bearer "+tokenStr)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)
	// shutdown server and db gracefully
	s.Close()

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expect 404 return code when trying to get an unexisting docker registry. "+
			"Got %d", resp.StatusCode)
	}
	var apiErr api.ApiError
	d := json.NewDecoder(resp.Body)
	d.Decode(&apiErr)
	if apiErr.Msg != expectedErrorMsg {
		t.Fatalf("Expect %s message when trying to get an unexisting docker registry. "+
			"Got %s", expectedErrorMsg, apiErr.Msg)
	}
}

// test creating a docker registry without authentication
func TestCreateANewDockerRegistryNotAuthenticated(t *testing.T) {
	// given
	t.SkipNow()

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	defer tests.CleanPersistence(conf)

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

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 401 {
		t.Fatalf("Expect 401 return code when trying to create a docker registry without auth. "+
			"Got %d", resp.StatusCode)
	}
}

// nominal test case for creating a docker registry
func TestCreateANewDockerRegistry(t *testing.T) {
	// given

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	defer tests.CleanPersistence(conf)

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	defer s.Close()

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
	if newRegistry.User != "tester" {
		t.Fatalf("expected User to have been perserved but got %s", newRegistry.User)
	}
	if newRegistry.Password != "" {
		t.Fatalf("expected Password to have been removed but got %s", newRegistry.Password)
	}
	if newRegistry.LastModificationTime == "" {
		t.Fatal("expected LastModificationTime to have initialized but is still 0")
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
