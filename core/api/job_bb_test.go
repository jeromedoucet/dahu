package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/dgrijalva/jwt-go"
	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/api"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/tests"
)

func TestCreateANewJobShouldReturn401WithoutAToken(t *testing.T) {
	// given
	job := model.Job{Name: "dahu", Url: "git@github.com:jeromedoucet/dahu.git", ImageName: "dahuci/dahu"}
	body, _ := json.Marshal(job)
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// when
	resp, err := http.Post(fmt.Sprintf("%s/jobs", s.URL),
		"application/json", bytes.NewBuffer(body))

	// shutdown server and db gracefully
	s.Close()
	tests.CleanPersistence(conf)

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
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// when
	resp, err := http.Get(fmt.Sprintf("%s/jobs", s.URL))

	// shutdown server and db gracefully
	s.Close()
	tests.CleanPersistence(conf)

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

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	job := model.Job{Name: "dahu", Url: "git@github.com:jeromedoucet/dahu.git", ImageName: "dahuci/dahu"}
	body, _ := json.Marshal(job)
	tokenStr := getToken("other_secret", time.Now().Add(1*time.Minute))
	req := buildJobsPostReq(body, tokenStr, s.URL)

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

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	tokenStr := getToken("other_secret", time.Now().Add(1*time.Minute))
	req := buildJobsGetReq(tokenStr, s.URL)

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

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	job := model.Job{Name: "dahu", Url: "git@github.com:jeromedoucet/dahu.git", ImageName: "dahuci/dahu"}
	body, _ := json.Marshal(job)
	tokenStr := getToken(conf.ApiConf.Secret, time.Now().Add(-1*time.Minute))
	req := buildJobsPostReq(body, tokenStr, s.URL)

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

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	tokenStr := getToken(conf.ApiConf.Secret, time.Now().Add(-1*time.Minute))
	req := buildJobsGetReq(tokenStr, s.URL)

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
		t.Fatalf("Expect 401 return code when trying to create a job with outdated credentials"+
			"without a token. Got %d", resp.StatusCode)
	}
}

func TestCreateANewJobShouldReturn400WhenNoName(t *testing.T) {
	// given

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	job := model.Job{Url: "git@github.com:jeromedoucet/dahu.git", ImageName: "dahuci/dahu"}
	body, _ := json.Marshal(job)
	tokenStr := getToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req := buildJobsPostReq(body, tokenStr, s.URL)

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
	if resp.StatusCode != 400 {
		t.Fatalf("Expect 400 return code when trying to create a job without name. "+
			"Got %d", resp.StatusCode)
	}
}

func TestCreateANewJobShouldReturn400WhenNoUrl(t *testing.T) {
	// given

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	job := model.Job{Name: "dahu", ImageName: "dahuci/dahu"}
	body, _ := json.Marshal(job)
	tokenStr := getToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req := buildJobsPostReq(body, tokenStr, s.URL)

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
	if resp.StatusCode != 400 {
		t.Fatalf("Expect 400 return code when trying to create a job without url. "+
			"Got %d", resp.StatusCode)
	}
}

func TestCreateANewJobShouldReturn400WhenNoImage(t *testing.T) {
	// given

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	job := model.Job{Url: "git@github.com:jeromedoucet/dahu.git", Name: "dahu"}
	body, _ := json.Marshal(job)
	tokenStr := getToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req := buildJobsPostReq(body, tokenStr, s.URL)

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

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// close the db for having an error
	close(conf.Close)

	// request setup
	job := model.Job{Name: "dahu", Url: "git@github.com:jeromedoucet/dahu.git", ImageName: "dahuci/dahu"}
	body, _ := json.Marshal(job)
	tokenStr := getToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req := buildJobsPostReq(body, tokenStr, s.URL)

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

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	job := model.Job{Name: "dahu", Url: "git@github.com:jeromedoucet/dahu.git", ImageName: "dahuci/dahu"}
	body, _ := json.Marshal(job)
	tokenStr := getToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req := buildJobsPostReq(body, tokenStr, s.URL)
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
		t.Fatalf("Expect 201 return code when trying to create a job. "+
			"Got %d", resp.StatusCode)
	}
	var dj model.Job
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&dj)
	if dj.Name != job.Name {
		t.Errorf("expected Name %s from file to equals %s", dj.Name, job.Name)
	}
	if dj.Url != job.Url {
		t.Errorf("expected Name %s from file to equals %s", dj.Url, job.Url)
	}
}

func TestListJobsShouldReturnAllJobs(t *testing.T) {
	// given

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	jobs := generateJobs(4)
	for _, job := range jobs {
		tests.InsertObject(conf, []byte("jobs"), []byte(job.Id), job)
	}

	// ap start
	s := httptest.NewServer(api.InitRoute(conf).Handler())

	// request setup
	tokenStr := getToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req := buildJobsGetReq(tokenStr, s.URL)
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
		if job.Url != j.Url {
			t.Fatalf("Expect to have a job with Url %s but got %s", job.Url, j.Url)
		}
		if job.ImageName != j.ImageName {
			t.Fatalf("Expect to have a job with ImageName %s but got %s", job.ImageName, j.ImageName)
		}
	}
}

func TestRunAJob(t *testing.T) {
	// given

	// configuration
	conf := configuration.InitConf()
	conf.ApiConf.Port = 4444
	conf.ApiConf.Secret = "secret"
	job := model.Job{Name: "dahu", Url: "git@github.com:jeromedoucet/dahu.git", ImageName: "dahuci/job-test"}
	job.EnvParam = make(map[string]string)
	job.EnvParam["STATUS"] = "success"
	job.GenerateId()
	tests.InsertObject(conf, []byte("jobs"), []byte(job.Id), job)
	a := api.InitRoute(conf)

	// ap start
	s := httptest.NewServer(a.Handler())

	// request setup
	reqBody := model.RunRequest{}
	body, _ := json.Marshal(reqBody)
	tokenStr := getToken(conf.ApiConf.Secret, time.Now().Add(1*time.Minute))
	req := buildJobTrigReq(body, tokenStr, s.URL, string(job.Id))
	cli := &http.Client{}

	// when
	resp, err := cli.Do(req)
	// shutdown server and db gracefully
	s.Close()
	a.Close()
	os.Remove(conf.PersistenceConf.Name)

	// then
	if err != nil {
		t.Fatalf("Expect to have to error, but got %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expect 200 return code when trying to strigger a Run. "+
			"Got %d", resp.StatusCode)
	}
}

func getToken(secret string, exp time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": exp.Unix(),
	})
	res, _ := token.SignedString([]byte(secret))
	return res
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

func buildJobTrigReq(body []byte, token, addr, jobId string) *http.Request {
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/jobs/%s/run",
		addr, jobId), bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)
	return req
}

func generateJobs(nbJobs int) []model.Job {
	jobs := make([]model.Job, nbJobs)
	for i := 0; i < nbJobs; i++ {
		jobs[i] = model.Job{Name: randomdata.SillyName(), Url: randomdata.IpV4Address(), ImageName: "someImage"}
		jobs[i].GenerateId()
	}
	return jobs
}

// todo test time out
