package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	bolt "github.com/coreos/bbolt"
	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
)

var singletonMutex = &sync.Mutex{}

var inMemorySingleton *inMemory

// return the inMemory singleton instance.
// create it if necessary
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
	dbInitialization(db)
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
		inMemorySingleton = nil
	}()
}

func dbInitialization(db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		return createBucketsIfNeeded(tx)
	})
}

// todo test me with errors case
func createBucketsIfNeeded(tx bucketCreationTransaction) error {
	var err error
	_, err = tx.CreateBucketIfNotExists([]byte("jobs"))
	if err != nil {
		return fmt.Errorf("ERROR >> job bucket creation failed : %s", err)
	}
	_, err = tx.CreateBucketIfNotExists([]byte("users"))
	if err != nil {
		return fmt.Errorf("ERROR >> user bucket creation failed : %s", err)
	}
	return nil
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
		// todo check that job is non-nil
		var updateErr error
		b := tx.Bucket([]byte("jobs"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing jobs. The database may be corrupted !")
		}
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
	if err == nil {
		return job, nil
	} else {
		return nil, err
	}

}

func (i *inMemory) CreateJobRun(jobRun *model.JobRun, jobId []byte, ctx context.Context) (*model.JobRun, error) {
	err := i.doUpdateAction(func(tx *bolt.Tx) error {
		var updateErr error
		if jobRun == nil || !jobRun.IsValid() {
			return errors.New(fmt.Sprintf("persistence >> trying to insert a nil or a an invalid JobRun : %+v", jobRun))
		}
		b := tx.Bucket([]byte("jobs"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing jobs. The database may be corrupted !")
		}
		var job *model.Job
		job, updateErr = getExistingJob(tx, jobId)
		if updateErr != nil { // todo test me
			return updateErr
		}
		updateErr = jobRun.GenerateId()
		if updateErr != nil {
			return updateErr
		}
		jobRun.Version = time.Now().UnixNano()
		job.AppendJobRun(jobRun)
		var data []byte
		data, updateErr = json.Marshal(job)
		if updateErr != nil {
			return updateErr
		}
		updateErr = b.Put(jobId, data) // todo handle the error
		return updateErr
	})
	if err == nil {
		return jobRun, nil
	} else {
		return nil, err
	}
}

func (i *inMemory) UpdateJobRun(jobRun *model.JobRun, jobId []byte, ctx context.Context) (*model.JobRun, error) {
	// if no more JobRun, return NoMorePersisted
	// if the JobRun si outdated, return the upToDateVersion and OutDated.
	// imagine what to do in case of outdated (depends what we were trying to update)
	err := i.doUpdateAction(func(tx *bolt.Tx) error {
		var updateErr error
		b := tx.Bucket([]byte("jobs"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing jobs. The database may be corrupted !")
		}
		if jobRun == nil || !jobRun.IsValid() { // todo test me
			return errors.New(fmt.Sprintf("persistence >> trying to insert a nil or a an invalid JobRun : %+v", jobRun))
		}
		var job *model.Job
		job, updateErr = getExistingJob(tx, jobId)
		if updateErr != nil {
			return updateErr
		}
		job.UpdateJobRun(jobRun) // todo find a way to return it
		var data []byte
		data, updateErr = json.Marshal(job)
		if updateErr != nil {
			return updateErr
		}
		updateErr = b.Put(jobId, data) // todo handle the error
		return updateErr
	})
	if err == nil {
		return jobRun, nil
	} else {
		return nil, err
	}
}

func getExistingJob(tx *bolt.Tx, jobId []byte) (*model.Job, error) {
	b := tx.Bucket([]byte("jobs"))
	if b == nil {
		return nil, errors.New("persistence >> CRITICAL error. No bucket for storing jobs. The database may be corrupted !")
	}
	var job model.Job
	jobData := b.Get(jobId)
	err := json.Unmarshal(jobData, &job) // todo handle this error
	return &job, err
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
	if err == nil {
		return &job, nil
	} else {
		return nil, err
	}
}

func (i *inMemory) GetJobs(ctx context.Context) ([]*model.Job, error) {
	// todo add missing tests
	jobs := make([]*model.Job, 0)
	err := i.doViewAction(func(tx *bolt.Tx) error {
		var mErr error
		b := tx.Bucket([]byte("jobs"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing jobs. The database may be corrupted !")
		}
		c := b.Cursor()
		jobs, mErr = doFetchJobs(c, jobs)
		return mErr
	})
	if err == nil {
		return jobs, nil
	} else {
		return nil, err
	}
}

func doFetchJobs(c *bolt.Cursor, jobs []*model.Job) ([]*model.Job, error) {
	res := jobs
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var job model.Job
		fmt.Println(fmt.Sprintf("persistence >> DEBUG. GetJobs: fetched value : %s, with key %s which is valid : %v", v, k, json.Valid(v)))
		mErr := json.Unmarshal(v, &job)
		if mErr != nil {
			return nil, mErr
		} else {
			res = append(res, &job)
		}
	}
	return res, nil
}

func (i *inMemory) GetUser(id string, ctx context.Context) (*model.User, error) {
	var user model.User
	err := i.doViewAction(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if b == nil {
			return errors.New("persistence >> CRITICAL error. No bucket for storing users. The database may be corrupted !")
		}
		data := b.Get([]byte(id))
		mErr := json.Unmarshal(data, &user)
		return mErr
	})
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
