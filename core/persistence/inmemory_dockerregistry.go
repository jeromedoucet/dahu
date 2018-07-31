package persistence

import (
	"context"

	"github.com/jeromedoucet/dahu/core/model"
)

func (i *inMemory) CreateDockerRegistry(job *model.DockerRegistry, ctx context.Context) (*model.DockerRegistry, error) {
	return nil, nil
}

func (i *inMemory) getDockerRegistry(id []byte, ctx context.Context) (*model.DockerRegistry, error) {
	return nil, nil
}

func (i *inMemory) getDockerRegistries(ctx context.Context) ([]*model.DockerRegistry, error) {
	return nil, nil
}

func (i *inMemory) deleteDockerRegistry(id []byte) error {
	return nil
}
