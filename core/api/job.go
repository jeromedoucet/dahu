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
	defer func() {
		if r := recover(); r != nil {
			log.Printf("WARN >> handleJobs encounter error : %+v", r)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	err := a.checkToken(r)
	if err == nil {
		var reqJob model.Job
		d := json.NewDecoder(r.Body)
		d.Decode(&reqJob)
		if !reqJob.IsValid() {
			log.Printf("WARN >> handleJobs encounter error : %+v is not valid", reqJob)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var newJob *model.Job
		newJob, err = a.repository.CreateJob(&reqJob, ctx)
		if err != nil {
			log.Printf("WARN >> handleJobs encounter error : %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		body, _ := json.Marshal(newJob) // todo handle err
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "%s", body)
		w.Write(body)
	} else {
		log.Printf("WARN >> handleJobs encounter error : %s ", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
	}
}
