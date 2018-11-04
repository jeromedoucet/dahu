package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	job_processing "github.com/jeromedoucet/dahu/core/job"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
	"github.com/jeromedoucet/route"
)

type execution struct {
	Branch string `json:"branch"`
}

type executionResult struct {
	Id string `json:"id"`
}

// handle request on jobs/
func (a *Api) handleJobs(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		a.onCreateJob(ctx, w, r)
	} else if r.Method == http.MethodGet {
		a.onGetJobs(ctx, w, r)
	} else {
		// todo return appropriate http code with a corresponding test
	}
}

func (a *Api) onCreateJob(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var reqJob model.Job
	var err error
	d := json.NewDecoder(r.Body)
	d.Decode(&reqJob)
	if !reqJob.IsValid() {
		log.Printf("ERROR >> createJob encounter error : %+v is not valid", reqJob)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var newJob *model.Job
	newJob, err = a.repository.CreateJob(&reqJob, ctx)
	if err != nil {
		log.Printf("ERROR >> createJob encounter error : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ := json.Marshal(newJob) // todo handle err
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (a *Api) onGetJobs(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	var jobs []*model.Job
	jobs, err = a.repository.GetJobs(ctx)
	if err != nil {
		log.Printf("ERROR >> GetJobs encounter error : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, job := range jobs {
		job.ToPublicModel()
	}
	body, err := json.Marshal(jobs)
	if err != nil {
		log.Printf("ERROR >> GetJobs encounter error : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (a *Api) onStartJob(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err persistence.PersistenceError
	var exec execution
	d := json.NewDecoder(r.Body)
	d.Decode(&exec)

	var job *model.Job
	path := route.SplitPath(r.URL.Path)
	jobId := path[len(path)-2]

	job, err = a.repository.GetJob([]byte(jobId), ctx)
	if err != nil {
		log.Printf("ERROR >> onStartJob encounter error : %s", err.Error())
		body := fromErrorToJson(err)
		if err.ErrorType() == persistence.NotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(body)
		return
	}

	log.Printf("INFO >> onStartJob asked for job id %s", string(job.Id))
	jobExecution := job_processing.Start(*job, exec.Branch, a.conf, ctx)
	log.Printf("INFO >> onStartJob start execution %s", jobExecution.Id)

	result := executionResult{Id: jobExecution.Id}

	body, marshErr := json.Marshal(result)
	if marshErr != nil {
		log.Printf("ERROR >> startJob encounter error : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (a *Api) onCancelJobExecution(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	path := route.SplitPath(r.URL.Path)
	jobId := path[len(path)-4]
	executionId := path[len(path)-2]
	log.Printf("INFO >> onCancelJobExecution asked for job id %s and job execution id : %s", jobId, executionId)

	job_processing.AskForCancelation(jobId, executionId)

	w.WriteHeader(http.StatusOK)
}

func (a *Api) onJobEventRegistration(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ws, err := a.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("ERROR >> onJobEventRegistration encounter error : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	path := route.SplitPath(r.URL.Path)
	jobId := path[len(path)-2]
	job_processing.AddWsEventListener(jobId, ws)
}
