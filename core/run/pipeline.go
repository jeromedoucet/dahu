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
	// todo choose the implementation witch fit the job
	pipeline := simplePipeline{repository: repository, job: job}
	return &pipeline
}

// fetch the code from a distant repo
// and return the name of the named volume
func fetchCode(imageName string, env map[string]string) string {
	// todo generate a named volume
	return ""
}
