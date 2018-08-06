package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/api"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/tests"
)

func TestCreateANewJobShouldReturn401WithoutAToken(t *testing.T) {
	// given
	sshAuth := model.SshAuthConfig{Url: "git@some-domain/some-repo.git", Key: "some-key", KeyPassword: "some-password"}
	scmConf := model.GitConfig{SshAuth: &sshAuth}
	job := model.Job{Name: "dahu", GitConf: scmConf}
	body, _ := json.Marshal(job)
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	defer s.Close()

	// when
	resp, err := http.Post(fmt.Sprintf("%s/jobs", s.URL),
		"application/json", bytes.NewBuffer(body))

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}

	if resp.StatusCode != 401 {
		t.Fatalf("Expect 401 return code when trying to create a job "+
			"without a token. Got %d", resp.StatusCode)
	}
}

func TestListJobsShouldReturn401WithoutAToken(t *testing.T) {
	// given
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	defer tests.CleanPersistence(conf)
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	defer s.Close()

	// when
	resp, err := http.Get(fmt.Sprintf("%s/jobs", s.URL))

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}

	if resp.StatusCode != 401 {
		t.Fatalf("Expect 401 return code when trying to create a job "+
			"without a token. Got %d", resp.StatusCode)
	}
}

func TestCreateANewJobShouldReturn401WhenBadCredentials(t *testing.T) {
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
	sshAuth := model.SshAuthConfig{Url: "git@some-domain/some-repo.git", Key: "some-key", KeyPassword: "some-password"}
	scmConf := model.GitConfig{SshAuth: &sshAuth}
	job := model.Job{Name: "dahu", GitConf: scmConf}
	body, _ := json.Marshal(job)
	tokenStr := tests.GetToken("other_secret", time.Now().Add(1*time.Minute))
	req := buildJobsPostReq(body, tokenStr, s.URL)

	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 401 {
		t.Fatalf("Expect 401 return code when trying to create a job with bad credentials"+
			"without a token. Got %d", resp.StatusCode)
	}
}

func TestListJobsShouldReturn401WhenBadCredentials(t *testing.T) {
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
	tokenStr := tests.GetToken("other_secret", time.Now().Add(1*time.Minute))
	req := buildJobsGetReq(tokenStr, s.URL)

	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 401 {
		t.Fatalf("Expect 401 return code when trying to create a job with bad credentials"+
			"without a token. Got %d", resp.StatusCode)
	}
}

func TestCreateANewJobShouldReturn401WhenTokenOutDated(t *testing.T) {
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
	sshAuth := model.SshAuthConfig{Url: "git@some-domain/some-repo.git", Key: "some-key", KeyPassword: "some-password"}
	scmConf := model.GitConfig{SshAuth: &sshAuth}
	job := model.Job{Name: "dahu", GitConf: scmConf}
	body, _ := json.Marshal(job)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(-1*time.Minute))
	req := buildJobsPostReq(body, tokenStr, s.URL)

	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 401 {
		t.Fatalf("Expect 401 return code when trying to create a job with outdated credentials"+
			"without a token. Got %d", resp.StatusCode)
	}
}

func TestListJobsShouldReturn401WhenTokenOutDated(t *testing.T) {
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
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(-1*time.Minute))
	req := buildJobsGetReq(tokenStr, s.URL)

	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 401 {
		t.Fatalf("Expect 401 return code when trying to create a job with outdated credentials"+
			"without a token. Got %d", resp.StatusCode)
	}
}

func TestCreateANewJobShouldReturn400WhenInvalidJob(t *testing.T) {
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
	sshAuth := model.SshAuthConfig{Url: "git@some-domain/some-repo.git", Key: "some-key", KeyPassword: "some-password"}
	scmConf := model.GitConfig{SshAuth: &sshAuth}
	job := model.Job{GitConf: scmConf}
	body, _ := json.Marshal(job)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req := buildJobsPostReq(body, tokenStr, s.URL)

	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 400 {
		t.Fatalf("Expect 400 return code when trying to create a job without name. "+
			"Got %d", resp.StatusCode)
	}
}

func TestCreateANewJobShouldReturn500WhenErroOnPersistenceLayer(t *testing.T) {
	// given

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	defer tests.DeletePersistence(conf)

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	defer s.Close()

	// close the db for having an error
	tests.ClosePersistence(conf)

	// request setup
	sshAuth := model.SshAuthConfig{Url: "git@some-domain/some-repo.git", Key: "some-key", KeyPassword: "some-password"}
	scmConf := model.GitConfig{SshAuth: &sshAuth}
	job := model.Job{Name: "dahu", GitConf: scmConf}
	body, _ := json.Marshal(job)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req := buildJobsPostReq(body, tokenStr, s.URL)

	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 500 {
		t.Fatalf("Expect 500 return code when error on persistence. "+
			"Got %d", resp.StatusCode)
	}
}

func TestCreateANewJobShouldCreateAndPersistAJob(t *testing.T) {
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
	sshAuth := model.SshAuthConfig{Url: "git@some-domain/some-repo.git", Key: "some-key", KeyPassword: "some-password"}
	scmConf := model.GitConfig{SshAuth: &sshAuth}
	job := model.Job{Name: "dahu", GitConf: scmConf}
	body, _ := json.Marshal(job)
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req := buildJobsPostReq(body, tokenStr, s.URL)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 201 {
		t.Fatalf("Expect 201 return code when trying to create a job. "+
			"Got %d", resp.StatusCode)
	}
	var dj model.Job
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&dj)
	if dj.Name != job.Name {
		t.Errorf("expected Name %s from file to equals %s", dj.Name, job.Name)
	}
}

func TestListJobsShouldReturnAllJobs(t *testing.T) {
	// given

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	defer tests.CleanPersistence(conf)
	jobs := generateJobs(4)
	for _, job := range jobs {
		tests.InsertObject(conf, []byte("jobs"), []byte(job.Id), job)
	}

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())
	defer s.Close()

	// request setup
	tokenStr := tests.GetToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req := buildJobsGetReq(tokenStr, s.URL)
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expect 200 return code when trying to list jobs."+
			"Got %d", resp.StatusCode)
	}
	var dj []model.Job
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&dj)
	mappedJob := make(map[string]model.Job)
	for _, job := range dj {
		mappedJob[string(job.Id)] = job
	}

	if len(dj) != len(jobs) {
		t.Fatalf("Expect to have %d jobs returned. Got %d", len(jobs), len(dj))
	}

	// verify the content of returned job
	for _, job := range jobs {
		j, _ := mappedJob[string(job.Id)]
		if string(job.Name) != string(j.Name) {
			t.Fatalf("Expect to have a job with Id %s but got %s", job.Name, j.Name)
		}
		if j.GitConf.SshAuth.Url != "git@some-domain/some-repo.git" {
			t.Fatalf("Expect to return ssh repo url in all jobs got %s for %s", j.GitConf.SshAuth.Url, j.Name)
		}
		if j.GitConf.SshAuth.Key != "" {
			t.Fatalf("Expect to return empty private key in all jobs got %s for %s", j.GitConf.SshAuth.Key, j.Name)
		}
		if j.GitConf.SshAuth.KeyPassword != "" {
			t.Fatalf("Expect to return empty private key password in all jobs got %s for %s", j.GitConf.SshAuth.KeyPassword, j.Name)
		}
	}
}

func buildJobsPostReq(body []byte, token string, addr string) *http.Request {
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/jobs",
		addr), bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)
	return req
}

func buildJobsGetReq(token string, addr string) *http.Request {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/jobs", addr), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	return req
}

func generateJobs(nbJobs int) []model.Job {
	jobs := make([]model.Job, nbJobs)
	sshAuth := model.SshAuthConfig{Url: "git@some-domain/some-repo.git", Key: "some-key", KeyPassword: "some-password"}
	scmConf := model.GitConfig{SshAuth: &sshAuth}
	for i := 0; i < nbJobs; i++ {
		jobs[i] = model.Job{Name: randomdata.SillyName(), GitConf: scmConf}
		jobs[i].GenerateId()
	}
	return jobs
}
