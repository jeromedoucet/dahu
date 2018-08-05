package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	bolt "github.com/coreos/bbolt"
	"github.com/jeromedoucet/dahu/core/model"
)

func (i *inMemory) CreateDockerRegistry(registry *model.DockerRegistry, ctx context.Context) (*model.DockerRegistry, PersistenceError) {
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
		return nil, wrapError(err)
	}
}

func (i *inMemory) GetDockerRegistry(id []byte, ctx context.Context) (*model.DockerRegistry, PersistenceError) {
	var registry model.DockerRegistry
	err := i.doViewAction(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("dockerRegistries"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing docker registries. The database may be corrupted !")
		}
		data := b.Get(id)
		if data == nil {
			return newPersistenceError(fmt.Sprintf("No docker registry with id %s found", string(id)), NotFound)
		}
		mErr := json.Unmarshal(data, &registry)
		return mErr
	})
	if err == nil {
		return &registry, nil
	} else {
		return nil, wrapError(err)
	}
}

func (i *inMemory) GetDockerRegistries(ctx context.Context) ([]*model.DockerRegistry, PersistenceError) {
	return nil, nil
}

func (i *inMemory) DeleteDockerRegistry(id []byte) PersistenceError {
	return nil
}
