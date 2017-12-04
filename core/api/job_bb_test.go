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
	var res model.JobRun
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&res)
	tests.RemoveContainer(res.ContainerName)
	if res.ContainerName == "" {
		t.Fatal("expect to have an container name but got ''")
	}
	if res.Status != model.RUNNING {
		t.Fatalf("exepct the JobRun status to be %d, but got %d", model.RUNNING, res.Status)
	}
	if len(res.Id) == 0 {
		t.Fatal("expect to have an id but got ''")
	}
	if res.StartTime == nil {
		t.Fatal("expect to have a start time, but got nil")
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

func buildJobTrigReq(body []byte, token, addr, jobId string) *http.Request {
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/jobs/%s/trigger",
		addr, jobId), bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)
	return req
}

// todo test time out
