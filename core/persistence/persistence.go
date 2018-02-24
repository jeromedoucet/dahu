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

	// job creation. If the job already has an id,
	// an error is returned.
	CreateJob(job *model.Job, ctx context.Context) (*model.Job, error)

	// get an existing job identified by the id parameter.
	GetJob(id []byte, ctx context.Context) (*model.Job, error)

	// get all existing jobs
	GetJobs(ctx context.Context) ([]*model.Job, error)

	// get an existing user identified by the id parameter.
	GetUser(id []byte, ctx context.Context) (*model.User, error)

	// will persist a jobRun on an existing Job.
	CreateJobRun(jobRun *model.JobRun, jobId []byte, ctx context.Context) (*model.JobRun, error)

	// will update a jobRun if it still exist on the Job
	UpdateJobRun(jobRun *model.JobRun, jobId []byte, ctx context.Context) (*model.JobRun, error)

	// this call will block until the underlying
	// connection or persistence system is open.
	WaitClose()
}

// return an instance of Repository configured with given configuration.
// for the moment, only the 'in-memory' persistence layer is
// supported
func GetRepository(conf *configuration.Conf) Repository {
	return getOrCreateInMemory(conf)
}
