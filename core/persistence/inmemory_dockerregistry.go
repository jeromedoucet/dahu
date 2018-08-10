package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

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
		// initialization of LastModificationDate field
		// that will be use later for optimistic lock on
		// update requests.
		registry.NewLastModificationTime()
		var data []byte
		data, updateErr = json.Marshal(registry)
		if updateErr != nil {
			return updateErr
		}
		updateErr = b.Put([]byte(registry.Id), data)
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

func (i *inMemory) DeleteDockerRegistry(id []byte) PersistenceError {
	err := i.doUpdateAction(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("dockerRegistries"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing docker registries. The database may be corrupted !")
		}
		// a get request is needed here because #Delete doesn't return an error
		// when key not found. This behavior is not consistent regarding the Api contract
		data := b.Get(id)
		if data == nil {
			return newPersistenceError(fmt.Sprintf("No docker registry with id %s found", string(id)), NotFound)
		}
		return b.Delete(id)
	})
	return wrapError(err)
}

func (i *inMemory) UpdateDockerRegistry(id []byte, registryUpdate *model.DockerRegistryUpdate, ctx context.Context) (*model.DockerRegistry, PersistenceError) {
	var updatedRegistry model.DockerRegistry
	err := i.doUpdateAction(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("dockerRegistries"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing docker registries. The database may be corrupted !")
		}
		var existingRegistry model.DockerRegistry
		data := b.Get(id)
		if data == nil {
			return newPersistenceError(fmt.Sprintf("No docker registry with id %s found", string(id)), NotFound)
		}
		mErr := json.Unmarshal(data, &existingRegistry)
		if mErr != nil {
			return mErr
		}
		// optimisitic lock check
		if existingRegistry.LastModificationTime != registryUpdate.LastModificationTime {
			updatedRegistry = existingRegistry
			return newPersistenceError(fmt.Sprintf("Conflict when trying to update registry with id %s", string(id)), Conflict)
		}
		registryUpdate.NewLastModificationTime()
		registry := registryUpdate.MergeForUpdate(&existingRegistry)
		data, updateErr := json.Marshal(registry)
		if updateErr != nil {
			return updateErr
		}
		updateErr = b.Put([]byte(registry.Id), data)
		updatedRegistry = *registry
		return updateErr
	})
	return &updatedRegistry, wrapError(err)
}

func (i *inMemory) GetDockerRegistries(ctx context.Context) ([]*model.DockerRegistry, PersistenceError) {
	registries := make([]*model.DockerRegistry, 0)
	err := i.doViewAction(func(tx *bolt.Tx) error {
		var mErr error
		b := tx.Bucket([]byte("dockerRegistries"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing docker registries. The database may be corrupted !")
		}
		c := b.Cursor()
		registries, mErr = doFetchDockerRegistries(c, registries)
		return mErr
	})
	if err == nil {
		return registries, nil
	} else {
		return nil, wrapError(err)
	}
	return nil, nil
}

func doFetchDockerRegistries(c *bolt.Cursor, registries []*model.DockerRegistry) ([]*model.DockerRegistry, error) {
	res := registries
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var registry model.DockerRegistry
		mErr := json.Unmarshal(v, &registry)
		if mErr != nil {
			return nil, mErr
		} else {
			res = append(res, &registry)
		}
	}
	// bbolt is not a relational db, so there is very little sorting features.
	// it is very simplier to do that directly in memory, considering the fact
	// that the size of read data is very small and the sort operation very simple
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res, nil
}
