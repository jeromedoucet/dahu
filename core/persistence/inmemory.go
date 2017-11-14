package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	bolt "github.com/coreos/bbolt"
	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
)

var singletonMutex = &sync.Mutex{}

var inMemorySingleton *inMemory

func getOrCreateInMemory(conf *configuration.Conf) Repository {
	// we may not use Once sync structure because the singleton
	// may be deleted and then recreated.
	// The sync strategy is agressive, but this
	// function may not be called very often, and there is
	// no other way to really prevents race conditions.
	singletonMutex.Lock()
	defer singletonMutex.Unlock()
	if inMemorySingleton == nil {
		createInMemory(conf)
	}
	return inMemorySingleton
}

// prepare a db inMemory instance
func createInMemory(conf *configuration.Conf) {
	inMemorySingleton = new(inMemory)
	inMemorySingleton.conf = conf
	inMemorySingleton.rwMutex = &sync.RWMutex{}
	db, _ := bolt.Open(conf.PersistenceConf.Name, 0600, nil) // todo handle this error
	db.Update(func(tx *bolt.Tx) error {
		var err error
		fmt.Println("test")
		_, err = tx.CreateBucketIfNotExists([]byte("jobs"))
		if err != nil {
			// todo test me
			return fmt.Errorf("create bucket: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			// todo test me
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	inMemorySingleton.db = db
	go func() {
		<-inMemorySingleton.conf.Close
		// when closing, the first thing
		// is to avoid any new call on #getOrCreateInMemory
		singletonMutex.Lock()
		defer singletonMutex.Unlock()
		// take a write lock. It permit to wait that all current
		// transaction finished
		inMemorySingleton.rwMutex.Lock()
		defer inMemorySingleton.rwMutex.Unlock()
		inMemorySingleton.db.Close() // todo handle this error
		// reset the singleton
		inMemorySingleton = nil
	}()
}

// In memory implementation
// of persistence.Repository.
// Use bbold as embedded db.
type inMemory struct {
	conf    *configuration.Conf
	db      *bolt.DB
	rwMutex *sync.RWMutex
}

func (i *inMemory) WaitClose() {
	<-i.conf.Close
}

func (i *inMemory) CreateJob(job *model.Job, ctx context.Context) (*model.Job, error) {
	err := i.doUpdateAction(func(tx *bolt.Tx) error {
		var updateErr error
		b := tx.Bucket([]byte("jobs"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing jobs. The database may be corrupted !")
		}
		// prepare & serialize the data
		updateErr = job.GenerateId()
		if updateErr != nil {
			return updateErr
		}
		var data []byte
		data, updateErr = json.Marshal(job)
		if updateErr != nil {
			return updateErr
		}
		updateErr = b.Put(job.Id, data)
		return updateErr
	})
	// if there is an error don't return the
	// job. This error must not be ignored
	if err == nil {
		return job, nil
	} else {
		return nil, err
	}

}

func (i *inMemory) CreateJobRun(jobRun *model.JobRun, jobId []byte, ctx context.Context) (*model.JobRun, error) {
	err := i.doUpdateAction(func(tx *bolt.Tx) error {
		var updateErr error
		b := tx.Bucket([]byte("jobs"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing jobs. The database may be corrupted !")
		}
		var job model.Job
		jobData := b.Get(jobId)
		updateErr = json.Unmarshal(jobData, &job) // todo handle this error
		// prepare & serialize the data
		updateErr = jobRun.GenerateId()
		if updateErr != nil {
			return updateErr
		}
		job.AppendJobRun(jobRun)
		var data []byte
		data, updateErr = json.Marshal(job)
		if updateErr != nil {
			return updateErr
		}
		updateErr = b.Put(jobId, data) // todo handle the error
		return updateErr
	})
	// if there is an error don't return the
	// job. This error must not be ignored
	if err == nil {
		return jobRun, nil
	} else {
		return nil, err
	}
}

func (i *inMemory) GetJob(id []byte, ctx context.Context) (*model.Job, error) {
	var job model.Job
	err := i.doViewAction(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("jobs"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing jobs. The database may be corrupted !")
		}
		data := b.Get(id)
		mErr := json.Unmarshal(data, &job)
		return mErr
	})
	// if there is an error don't return the
	// job. This error must not be ignored
	if err == nil {
		return &job, nil
	} else {
		return nil, err
	}
}

func (i *inMemory) GetUser(id []byte, ctx context.Context) (*model.User, error) {
	var user model.User
	err := i.doViewAction(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing users. The database may be corrupted !")
		}
		data := b.Get(id)
		mErr := json.Unmarshal(data, &user)
		return mErr
	})
	// if there is an error don't return the
	// job. This error must not be ignored
	if err == nil {
		return &user, nil
	} else {
		return nil, err
	}
}

// doUpdateAction will ensure that the
// db is available and then execute the function
// inside a read/write transaction. An error is returned
// if an issue appears.
func (i *inMemory) doUpdateAction(action func(tx *bolt.Tx) error) error {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()
	select {
	case <-i.conf.Close:
		return errors.New("persistence >> the database is close or closing. Operation impossible.")
	default:
		// only one read/write transaction is allowed.
		return i.db.Update(action)
	}

}

// doViewAction will ensure that the
// db is available and then execute the function
// inside a read transaction. An error is returned
// if an issue appears.
func (i *inMemory) doViewAction(action func(tx *bolt.Tx) error) error {
	i.rwMutex.RLock()
	defer i.rwMutex.RUnlock()
	select {
	case <-i.conf.Close:
		return errors.New("persistence >> the database is close or closing. Operation impossible.")
	default:
		return i.db.View(action)
	}

}
