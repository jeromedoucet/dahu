package run

import (
	"context"

	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
)

type Pipeline interface {
	Start(context.Context) error
	WaitUntilFinished()
}

func NewPipeline(job *model.Job, repository persistence.Repository) Pipeline {
	// todo choose the implementation witch firt the job
	pipeline := simplePipeline{repository: repository, job: job}
	return &pipeline
}
