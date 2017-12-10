package run

import (
	"context"
	"os"
	"time"

	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
)

type simplePipeline struct {
	repository     persistence.Repository
	job            *model.Job
	currentProcess *Process
}

func (pipeline *simplePipeline) Start(ctx context.Context) error {
	params := newProcessParams(pipeline.job)
	if params.Env == nil { // todo test cover me
		params.Env = make(map[string]string)
	}
	params.Env["REPO_URL"] = pipeline.job.Url // todo test cover me
	process := NewProcess(params, pipeline.repository)
	_, err := process.Start(ctx) // todo cover test for error
	pipeline.currentProcess = process
	return err
}

func (pipeline *simplePipeline) WaitUntilFinished() {
	// todo test me
	if pipeline.currentProcess != nil {
		<-pipeline.currentProcess.Done()
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
