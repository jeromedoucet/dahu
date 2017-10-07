package persistence

import (
	"context"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
)

// data access layer. Expose
// functions that allow manipulating data
// regardless of the underlying persistence system
type Repository interface {
	CreateJob(job *model.Job, ctx context.Context) (*model.Job, error)
	GetJob(id string, ctx context.Context)
	WaitClose()
}

// return an instance of Repository configured with given configuration.
// for the moment, only the 'in-memory' persistence layer is
// supported
func GetRepository(conf *configuration.Conf) Repository {
	return getOrCreateInMemory(conf)
}
