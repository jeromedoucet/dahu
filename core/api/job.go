package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jeromedoucet/dahu/core/model"
)

// handle request on jobs/
func (a *Api) handleJobs(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// recover from any panic coming form /jobs requests
	defer func() {
		if r := recover(); r != nil {
			log.Printf("WARN >> handleJobs encounter error : %+v", r)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	err := a.checkToken(r)
	if err == nil {
		if r.Method == http.MethodPost {
			a.onCreateJob(ctx, w, r)
		} else if r.Method == http.MethodGet {
			a.onGetJobs(ctx, w, r)
		} else {
			// todo return appropriate http code with a corresponding test
		}
	} else {
		log.Printf("WARN >> handleJobs encounter error : %s ", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
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
	// todo add tests
	var err error
	var jobs []*model.Job
	jobs, err = a.repository.GetJobs(ctx)
	if err != nil {
		log.Printf("ERROR >> GetJobs encounter error : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ := json.Marshal(jobs) // todo handle err
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", body)
	w.Write(body)
}
