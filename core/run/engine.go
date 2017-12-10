package run

import (
	"context"
	"errors"
	"sync"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
)

type RunEngine interface {
	StartOneJob(job *model.Job, ctx context.Context) error
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
func (re *SimpleRunEngine) StartOneJob(job *model.Job, ctx context.Context) error {
	re.runningCount.RLock()
	select {
	case <-re.conf.Close:
		re.runningCount.RUnlock()
		return errors.New("run >> the application is shutting down. Operation impossible.")
	default:
		pipeline := NewPipeline(job, re.repository)
		err := pipeline.Start(ctx)
		if err == nil {
			go func() {
				defer re.runningCount.RUnlock()
				pipeline.WaitUntilFinished()
			}()
		} else {
			re.runningCount.RUnlock()
		}
		return err
	}
}
