package run

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
)

type RunEngine interface {
	StartOneRun(job *model.Job, ctx context.Context) (*model.JobRun, error)
	WaitClose()
}

// create a new run engine.
func NewRunEngine(conf *configuration.Conf) RunEngine {
	r := new(SimpleRunEngine)
	r.runningCount = &sync.RWMutex{}
	r.conf = conf
	r.repository = persistence.GetRepository(conf)
	return r
}

// unique point of thruth
// for Run.
type SimpleRunEngine struct {
	conf         *configuration.Conf
	runningCount *sync.RWMutex
	repository   persistence.Repository
}

// Wait for all current run to be finished and then
// close the Close channel
func (r *SimpleRunEngine) WaitClose() {
	// todo test me
	// todo add a context
	r.runningCount.Lock()
	defer r.runningCount.Unlock()
	close(r.conf.Close)
}

// Start one new Run from a given job
func (re *SimpleRunEngine) StartOneRun(job *model.Job, ctx context.Context) (*model.JobRun, error) {
	re.runningCount.RLock()
	select {
	case <-re.conf.Close:
		re.runningCount.RUnlock()
		return nil, errors.New("run >> the application is shutting down. Operation impossible.")
	default:
		params := newProcessParams(job)
		if params.Env == nil { // todo test cover me
			params.Env = make(map[string]string)
		}
		params.Env["REPO_URL"] = job.Url // todo test cover me
		r := NewProcess(params, re.repository)
		res, err := r.Start(ctx) // todo cover test for error
		if err == nil {
			// if the run has started without error
			// defer the unlock in another goroutine
			go func() {
				defer re.runningCount.RUnlock()
				<-r.Done()
			}()
		} else {
			re.runningCount.RUnlock()
		}
		return &res, err
	}
}

// create a new Process Params from a given Job
func newProcessParams(job *model.Job) ProcessParams {
	// todo 1 generate an Id value
	// todo 2 outpout writer ?
	// todo 3 time out ?
	return ProcessParams{
		Image:        job.ImageName,
		Env:          job.EnvParam,
		OutputWriter: os.Stdout,
		TimeOut:      time.Second * 1,
		JobId:        job.Id,
	}

}
