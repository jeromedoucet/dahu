package api_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	tests_container "github.com/jeromedoucet/dahu-tests/container"
	"github.com/jeromedoucet/dahu-tests/ssh"
	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/api"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/tests"
)

func TestJobExecutionSuite(t *testing.T) {
	registryAuth := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "tester",
		Password: "test",
	}
	data, _ := json.Marshal(registryAuth)
	registryAuthString := base64.StdEncoding.EncodeToString(data)

	tests_container.PushImage(configuration.DockerApiVersion, "debian", "localhost:5000/debian", registryAuthString)
	beforeEach := func() {
		tests.CleanPersistence(configuration.InitConf())
		conf = configuration.InitConf()
		conf.ApiConf.Port = 4444
		conf.ApiConf.Secret = "secret"
	}

	afterEach := func() {
		tests.CleanPersistence(conf)
		if s != nil {
			s.Close()
		}
	}

	beforeEach()
	t.Run("job not found", jobNotFound)
	afterEach()

	beforeEach()
	t.Run("simple job with two step sucess execution", simpleJobSuccess)
	afterEach()

	beforeEach()
	t.Run("job with one step and service sucess execution", serviceJobSuccess)
	afterEach()

	beforeEach()
	t.Run("simple job with two step sucess execution and private registry", simpleJobSuccessWithPrivateRegistry)
	afterEach()

	beforeEach()
	t.Run("simple job failure on a step", simpleJobFailure)
	afterEach()

	beforeEach()
	t.Run("simple job failure on fetch", simpleJobFetchFailure)
	afterEach()

	beforeEach()
	t.Run("simple job cancelation", simpleJobCancel)
	afterEach()

	// TODO make sure the volume is not deleted => delete it at the end of the test
	// TODO check that all container are deleted
	// TODO clean of the listner content
	// TODO graceful stop of notifier
	// TODO test on another branch
	// TODO improve job execution saving
	// TODO enpoint to job execution => list, details, delete
}

func jobNotFound(t *testing.T) {
	// given
	// Start the server AFTER inserting the data.
	s = httptest.NewServer(api.InitRoute(conf).Handler())

	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	h := http.Header{}
	h.Add("Authorization", "Bearer "+tokenStr)

	execution := struct {
		Branch string `json:"branch"`
	}{
		"master",
	}
	body, _ := json.Marshal(execution)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/jobs/%s/executions", s.URL, "unknown"), bytes.NewBuffer(body))
	req.Header = h
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 404 {
		t.Fatalf("Expect 404 return code .Got %d", resp.StatusCode)
	}
}

// test a simple job execution with success
func simpleJobSuccess(t *testing.T) {
	// given
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/test-repo.git", gitRepoIp), Key: ssh.PrivateProtected, KeyPassword: "tester"}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	job := model.Job{
		Name:            "test",
		GitConf:         gitConfig,
		RemoveWorkspace: true,
		Steps: []model.Step{
			model.Step{
				Name:          "create file on debian",
				Image:         model.Image{Name: "debian"},
				Envs:          map[string]string{"HELLO": "hello world !"},
				Command:       []string{"/bin/sh", "-c", `echo "$HELLO" > build.out`},
				MountingPoint: "/build",
			},
			model.Step{
				Name:          "read file on fedora",
				Image:         model.Image{Name: "fedora"},
				Command:       []string{"cat", "build.out"},
				MountingPoint: "/build",
			},
		},
	}
	job.GenerateId()
	tests.InsertObject(conf, []byte("jobs"), job.Id, job)

	// Start the server AFTER inserting the data.
	s = httptest.NewServer(api.InitRoute(conf).Handler())

	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	h := http.Header{}
	h.Add("Authorization", "Bearer "+tokenStr)

	wsConn := openWsConn(s.URL, fmt.Sprintf("/jobs/%s/live", string(job.Id)), h, t)
	eventsChan := collectWsEvent(wsConn, t)

	execution := struct {
		Branch string `json:"branch"`
	}{
		"master",
	}
	body, _ := json.Marshal(execution)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/jobs/%s/executions", s.URL, job.Id), bytes.NewBuffer(body))
	req.Header = h
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expect 200 return code .Got %d", resp.StatusCode)
	}

	events := <-eventsChan
	if len(events) != 14 {
		t.Fatalf("expect %d events, got %d", 14, len(events))
	}

	if events[0].Type != model.JobStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 1, model.JobStart, events[0].Type)
	}

	if events[1].Type != model.StepStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 2, model.StepStart, events[1].Type)
	}

	if events[7].Type != model.StepSucceed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 8, model.StepSucceed, events[7].Type)
	}

	if events[8].Type != model.StepStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 9, model.StepStart, events[8].Type)
	}

	if events[9].Type != model.StepSucceed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 10, model.StepSucceed, events[9].Type)
	}

	if events[10].Type != model.StepStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 11, model.StepStart, events[10].Type)
	}

	if events[11].Value != "hello world !" {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 11, "hello world !", events[11].Value)
	}

	if events[12].Type != model.StepSucceed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 13, model.StepSucceed, events[12].Type)
	}

	if events[13].Type != model.JobSucceed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 14, model.JobStart, events[13].Type)
	}
}

// this case will check that a service (here nginx) may be
// available for one job.
func serviceJobSuccess(t *testing.T) {
	// given
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/test-repo.git", gitRepoIp), Key: ssh.PrivateProtected, KeyPassword: "tester"}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	job := model.Job{
		Name:            "test",
		GitConf:         gitConfig,
		RemoveWorkspace: true,
		Steps: []model.Step{
			model.Step{
				Name:          "create file on debian",
				Image:         model.Image{Name: "fedora"},
				Command:       []string{"/bin/sh", "-c", `curl -s -o /dev/null -w "%{http_code}\n" http://"${Nginx_Test_HOST}"`},
				MountingPoint: "/build",
				Services: []*model.Service{
					&model.Service{
						Name:  "Nginx_Test",
						Image: model.Image{Name: "nginx"},
						ExposedPorts: []*model.Port{
							&model.Port{Num: 80, Prototype: "http"},
						},
					},
				},
			},
		},
	}
	// "%{http_code}\n" http://"${Nginx_Test_HOST}"
	job.GenerateId()
	tests.InsertObject(conf, []byte("jobs"), job.Id, job)

	// Start the server AFTER inserting the data.
	s = httptest.NewServer(api.InitRoute(conf).Handler())

	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	h := http.Header{}
	h.Add("Authorization", "Bearer "+tokenStr)

	wsConn := openWsConn(s.URL, fmt.Sprintf("/jobs/%s/live", string(job.Id)), h, t)
	eventsChan := collectWsEvent(wsConn, t)

	execution := struct {
		Branch string `json:"branch"`
	}{
		"master",
	}
	body, _ := json.Marshal(execution)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/jobs/%s/executions", s.URL, job.Id), bytes.NewBuffer(body))
	req.Header = h
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expect 200 return code .Got %d", resp.StatusCode)
	}

	events := <-eventsChan
	if len(events) != 12 {
		t.Fatalf("expect %d events, got %d", 12, len(events))
	}

	if events[0].Type != model.JobStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 1, model.JobStart, events[0].Type)
	}

	if events[1].Type != model.StepStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 2, model.StepStart, events[1].Type)
	}

	if events[7].Type != model.StepSucceed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 8, model.StepSucceed, events[7].Type)
	}

	if events[8].Type != model.StepStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 8, model.StepStart, events[8].Type)
	}

	if events[9].Value != "200" {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 9, "200", events[9].Value)
	}

	if events[10].Type != model.StepSucceed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 10, model.StepSucceed, events[10].Type)
	}

	if events[11].Type != model.JobSucceed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 11, model.JobStart, events[11].Type)
	}
}

// same test than success, but with
// one image in a private registry
func simpleJobSuccessWithPrivateRegistry(t *testing.T) {
	// given
	registry := new(model.DockerRegistry)
	registry.Name = "test"
	registry.Url = "localhost:5000" // TODO use the docker network ip instead
	registry.User = "tester"
	registry.Password = "test"
	registry.GenerateId()
	registry.NewLastModificationTime()

	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/test-repo.git", gitRepoIp), Key: ssh.PrivateProtected, KeyPassword: "tester"}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	job := model.Job{
		Name:            "test",
		GitConf:         gitConfig,
		RemoveWorkspace: true,
		Steps: []model.Step{
			model.Step{
				Name:          "create file on debian",
				Image:         model.Image{Name: "debian", RegistryId: registry.Id},
				Command:       []string{"/bin/sh", "-c", "echo hello world > build.out"},
				MountingPoint: "/build",
			},
			model.Step{
				Name:          "read file on fedora",
				Image:         model.Image{Name: "fedora"},
				Command:       []string{"cat", "build.out"},
				MountingPoint: "/build",
			},
		},
	}
	job.GenerateId()

	tests.InsertObject(conf, []byte("dockerRegistries"), []byte(registry.Id), registry)
	tests.InsertObject(conf, []byte("jobs"), job.Id, job)

	// Start the server AFTER inserting the data.
	s = httptest.NewServer(api.InitRoute(conf).Handler())

	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	h := http.Header{}
	h.Add("Authorization", "Bearer "+tokenStr)

	wsConn := openWsConn(s.URL, fmt.Sprintf("/jobs/%s/live", string(job.Id)), h, t)
	eventsChan := collectWsEvent(wsConn, t)

	execution := struct {
		Branch string `json:"branch"`
	}{
		"master",
	}
	body, _ := json.Marshal(execution)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/jobs/%s/executions", s.URL, job.Id), bytes.NewBuffer(body))
	req.Header = h
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expect 200 return code .Got %d", resp.StatusCode)
	}

	events := <-eventsChan
	if len(events) != 14 {
		t.Fatalf("expect %d events, got %d", 14, len(events))
	}

	if events[0].Type != model.JobStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 1, model.JobStart, events[0].Type)
	}

	if events[1].Type != model.StepStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 2, model.StepStart, events[1].Type)
	}

	if events[7].Type != model.StepSucceed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 8, model.StepSucceed, events[7].Type)
	}

	if events[8].Type != model.StepStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 9, model.StepStart, events[8].Type)
	}

	if events[9].Type != model.StepSucceed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 10, model.StepSucceed, events[9].Type)
	}

	if events[10].Type != model.StepStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 11, model.StepStart, events[10].Type)
	}

	if events[12].Type != model.StepSucceed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 13, model.StepSucceed, events[12].Type)
	}

	if events[13].Type != model.JobSucceed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 14, model.JobStart, events[13].Type)
	}
}

// test a simple job execution with one step failure
func simpleJobFailure(t *testing.T) {
	// given
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/test-repo.git", gitRepoIp), Key: ssh.PrivateProtected, KeyPassword: "tester"}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	job := model.Job{
		Name:            "test",
		GitConf:         gitConfig,
		RemoveWorkspace: true,
		Steps: []model.Step{
			model.Step{
				Name:          "create file on debian",
				Image:         model.Image{Name: "debian"},
				Command:       []string{"/bin/sh", "-c", "this-is-a-bad-command"},
				MountingPoint: "/build",
			},
			model.Step{
				Name:          "read file on fedora",
				Image:         model.Image{Name: "fedora"},
				Command:       []string{"cat", "build.out"},
				MountingPoint: "/build",
			},
		},
	}
	job.GenerateId()
	tests.InsertObject(conf, []byte("jobs"), job.Id, job)

	// Start the server AFTER inserting the data.
	s = httptest.NewServer(api.InitRoute(conf).Handler())

	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	h := http.Header{}
	h.Add("Authorization", "Bearer "+tokenStr)

	wsConn := openWsConn(s.URL, fmt.Sprintf("/jobs/%s/live", string(job.Id)), h, t)
	eventsChan := collectWsEvent(wsConn, t)

	execution := struct {
		Branch string `json:"branch"`
	}{
		"master",
	}
	body, _ := json.Marshal(execution)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/jobs/%s/executions", s.URL, job.Id), bytes.NewBuffer(body))
	req.Header = h
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expect 200 return code .Got %d", resp.StatusCode)
	}

	events := <-eventsChan
	if len(events) != 12 {
		t.Fatalf("expect %d events, got %d", 12, len(events))
	}

	if events[0].Type != model.JobStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 1, model.JobStart, events[0].Type)
	}

	if events[1].Type != model.StepStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 2, model.StepStart, events[1].Type)
	}

	if events[7].Type != model.StepSucceed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 8, model.StepSucceed, events[7].Type)
	}

	if events[8].Type != model.StepStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 9, model.StepStart, events[8].Type)
	}

	if events[10].Type != model.StepFailed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 11, model.StepFailed, events[10].Type)
	}

	if events[11].Type != model.JobFailed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 12, model.StepFailed, events[11].Type)
	}

}

// test a simple job execution with one step failure
func simpleJobFetchFailure(t *testing.T) {
	// given
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/unknown-repo.git", gitRepoIp), Key: ssh.PrivateProtected, KeyPassword: "tester"}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	job := model.Job{
		Name:            "test",
		GitConf:         gitConfig,
		RemoveWorkspace: true,
		Steps: []model.Step{
			model.Step{
				Name:          "create file on debian",
				Image:         model.Image{Name: "debian"},
				Command:       []string{"/bin/sh", "-c", "echo hello world > build.out"},
				MountingPoint: "/build",
			},
			model.Step{
				Name:          "read file on fedora",
				Image:         model.Image{Name: "fedora"},
				Command:       []string{"cat", "build.out"},
				MountingPoint: "/build",
			},
		},
	}
	job.GenerateId()
	tests.InsertObject(conf, []byte("jobs"), job.Id, job)

	// Start the server AFTER inserting the data.
	s = httptest.NewServer(api.InitRoute(conf).Handler())

	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	h := http.Header{}
	h.Add("Authorization", "Bearer "+tokenStr)

	wsConn := openWsConn(s.URL, fmt.Sprintf("/jobs/%s/live", string(job.Id)), h, t)
	eventsChan := collectWsEvent(wsConn, t)

	execution := struct {
		Branch string `json:"branch"`
	}{
		"master",
	}
	body, _ := json.Marshal(execution)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/jobs/%s/executions", s.URL, job.Id), bytes.NewBuffer(body))
	req.Header = h
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expect 200 return code .Got %d", resp.StatusCode)
	}

	events := <-eventsChan
	if len(events) != 7 {
		t.Fatalf("expect %d events, got %d", 7, len(events))
	}

	if events[0].Type != model.JobStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 1, model.JobStart, events[0].Type)
	}

	if events[1].Type != model.StepStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 2, model.StepStart, events[1].Type)
	}

	if events[5].Type != model.StepFailed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 6, model.StepSucceed, events[5].Type)
	}

	if events[6].Type != model.JobFailed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 7, model.StepFailed, events[6].Type)
	}
}

// test a simple job execution cancelation
func simpleJobCancel(t *testing.T) {
	// given
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/test-repo.git", gitRepoIp), Key: ssh.PrivateProtected, KeyPassword: "tester"}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	job := model.Job{
		Name:            "test",
		GitConf:         gitConfig,
		RemoveWorkspace: true,
		Steps: []model.Step{
			model.Step{
				Name:          "create file on debian",
				Image:         model.Image{Name: "debian"},
				Command:       []string{"sleep", "infinity"},
				MountingPoint: "/build",
			},
			model.Step{
				Name:          "read file on fedora",
				Image:         model.Image{Name: "fedora"},
				Command:       []string{"cat", "build.out"},
				MountingPoint: "/build",
			},
		},
	}
	job.GenerateId()
	tests.InsertObject(conf, []byte("jobs"), job.Id, job)

	// Start the server AFTER inserting the data.
	s = httptest.NewServer(api.InitRoute(conf).Handler())

	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	h := http.Header{}
	h.Add("Authorization", "Bearer "+tokenStr)

	wsConn1 := openWsConn(s.URL, fmt.Sprintf("/jobs/%s/live", string(job.Id)), h, t)
	wsConn2 := openWsConn(s.URL, fmt.Sprintf("/jobs/%s/live", string(job.Id)), h, t)
	eventsChan := collectWsEvent(wsConn1, t)
	notifChan := wsEventNotification(wsConn2, model.StepStart, 2, t)

	execution := struct {
		Branch string `json:"branch"`
	}{
		"master",
	}
	body, _ := json.Marshal(execution)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/jobs/%s/executions", s.URL, job.Id), bytes.NewBuffer(body))
	req.Header = h
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expect 200 return code .Got %d", resp.StatusCode)
	}

	var executionResult struct {
		Id string `json:"id"`
	}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&executionResult)

	<-notifChan

	// cancel the job
	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/jobs/%s/executions/%s/cancelation", s.URL, job.Id, executionResult.Id), nil)
	req.Header = h
	cli = &http.Client{}

	// when
	resp, err = cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expect 200 return code .Got %d", resp.StatusCode)
	}

	// wait for all event to arrive
	events := <-eventsChan
	if len(events) != 11 {
		t.Fatalf("expect %d events, got %d", 11, len(events))
	}

	if events[0].Type != model.JobStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 1, model.JobStart, events[0].Type)
	}

	if events[1].Type != model.StepStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 2, model.StepStart, events[1].Type)
	}

	if events[7].Type != model.StepSucceed {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 8, model.StepSucceed, events[7].Type)
	}

	if events[8].Type != model.StepStart {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 9, model.StepStart, events[8].Type)
	}

	if events[9].Type != model.StepCanceled {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 10, model.StepCanceled, events[9].Type)
	}

	if events[10].Type != model.JobCanceled {
		t.Fatalf("expect event n %d, to be of type %s, but is %s", 11, model.JobCanceled, events[10].Type)
	}
}

// openWsConn allow to get a websocket connection
// for testing purposed
func openWsConn(serverUrl, path string, h http.Header, t *testing.T) *websocket.Conn {
	surl, _ := url.Parse(serverUrl)
	u := url.URL{Scheme: "ws", Host: surl.Hostname() + ":" + surl.Port(), Path: path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), h)
	if err != nil {
		t.Fatal(err.Error())
		return nil
	}
	return c
}

// collectWsEvent will collect job event till a
// terminal job event is received. This function
// is blocking.
func collectWsEvent(wsConn *websocket.Conn, t *testing.T) chan []model.Event {
	var events []model.Event
	c := make(chan []model.Event)
	go func() {
		for {
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				t.Fatal(err.Error())
				break
			}
			var e model.Event
			err = json.Unmarshal(message, &e)
			if err != nil {
				t.Fatal(err.Error())
				break
			}
			t.Logf("TEST >> receive event : %s - %s", e.Type, e.Value)
			events = append(events, e)
			if e.Type == model.JobSucceed || e.Type == model.JobFailed || e.Type == model.JobCanceled {
				break
			}
		}
		c <- events
		close(c)
	}()
	return c
}

// wsEventNotification will notify the caller
// when a specific event occurs for the n times
func wsEventNotification(wsConn *websocket.Conn, messageType model.EventType, nbMessage int, t *testing.T) chan interface{} {
	c := make(chan interface{})
	count := 0
	go func() {
		for {
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				t.Fatal(err.Error())
				break
			}
			var e model.Event
			err = json.Unmarshal(message, &e)
			if err != nil {
				t.Fatal(err.Error())
				break
			}
			if e.Type == messageType {
				count += 1
			}
			if count >= nbMessage {
				break
			}
		}
		c <- nil
		close(c)
	}()

	return c
}
