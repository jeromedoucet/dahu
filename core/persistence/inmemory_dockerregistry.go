package persistence

import (
	"context"
	"encoding/json"
	"errors"

	bolt "github.com/coreos/bbolt"
	"github.com/jeromedoucet/dahu/core/model"
)

func (i *inMemory) CreateDockerRegistry(registry *model.DockerRegistry, ctx context.Context) (*model.DockerRegistry, error) {
	// todo see if factorization can be done
	err := i.doUpdateAction(func(tx *bolt.Tx) error {
		// todo check that docker registry is non-nil
		var updateErr error
		b := tx.Bucket([]byte("dockerRegistries"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing docker registries. The database may be corrupted !")
		}
		updateErr = registry.GenerateId()
		if updateErr != nil {
			return updateErr
		}
		var data []byte
		data, updateErr = json.Marshal(registry)
		if updateErr != nil {
			return updateErr
		}
		updateErr = b.Put(registry.Id, data)
		return updateErr
	})
	if err == nil {
		return registry, nil
	} else {
		return nil, err
	}
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
