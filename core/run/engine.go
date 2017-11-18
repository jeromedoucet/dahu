package run

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
)

// create a new run engine.
func NewRunEngine(conf *configuration.Conf) *RunEngine {
	r := new(RunEngine)
	r.runningCount = &sync.RWMutex{}
	r.conf = conf
	return r
}

// unique point of thruth
// for Run.
type RunEngine struct {
	conf         *configuration.Conf
	runningCount *sync.RWMutex
}

// Wait for all current run to be finished and then
// close the Close channel
func (r *RunEngine) WaitClose() {
	// todo test me
	// todo add a context
	r.runningCount.Lock()
	defer r.runningCount.Unlock()
	close(r.conf.Close)
}

// Start one new Run from a given
// job
func (re *RunEngine) StartOneRun(job *model.Job, ctx context.Context) (model.JobRun, error) {
	re.runningCount.RLock()
	select {
	case <-re.conf.Close:
		re.runningCount.RUnlock()
		return model.JobRun{}, errors.New("run >> the application is shutting down. Operation impossible.")
	default:
		// todo 1 generate an Id value
		// todo 2 outpout writer ?
		// todo 3 time out ?
		params := ProcessParams{
			Id:           "test-2",
			Image:        job.ImageName,
			Env:          job.EnvParam,
			OutputWriter: os.Stdout,
			TimeOut:      time.Second * 1,
		}
		if params.Env == nil { // todo test cover me
			params.Env = make(map[string]string)
		}
		params.Env["REPO_URL"] = job.Url // todo test cover me
		r := NewProcess(params)
		res := model.JobRun{ContainerName: params.ContainerName()} // todo set run id !
		err := r.Start(ctx)                                        // todo cover test for error
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
		return res, err
	}
}
