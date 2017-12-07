package run

import (
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
)

type Pipeline interface {
}

func NewPipeline(job *model.Job, repository persistence.Repository) Pipeline {
	return nil
}
