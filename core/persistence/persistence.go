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
	// an PersistenceError is returned.
	CreateJob(job *model.Job, ctx context.Context) (*model.Job, PersistenceError)

	// get an existing job identified by the id parameter.
	GetJob(id []byte, ctx context.Context) (*model.Job, PersistenceError)

	// get all existing jobs
	GetJobs(ctx context.Context) ([]*model.Job, PersistenceError)

	// get an existing user identified by the id parameter.
	GetUser(id string, ctx context.Context) (*model.User, PersistenceError)

	// docker registry creation. If the docker regitry already has an id,
	// an PersistenceError is returned.
	CreateDockerRegistry(registry *model.DockerRegistry, ctx context.Context) (*model.DockerRegistry, PersistenceError)

	// get an existing docker registry identified by the id parameter.
	GetDockerRegistry(id []byte, ctx context.Context) (*model.DockerRegistry, PersistenceError)

	// get all existing docker registries.
	GetDockerRegistries(ctx context.Context) ([]*model.DockerRegistry, PersistenceError)

	// delete one existing docker registry
	DeleteDockerRegistry(id []byte) PersistenceError

	// update one existing docker registry
	UpdateDockerRegistry(id []byte, registry *model.DockerRegistryUpdate, ctx context.Context) (*model.DockerRegistry, PersistenceError)

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
