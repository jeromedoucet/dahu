package persistence

import (
	"context"
	"encoding/json"
	"errors"

	bolt "github.com/coreos/bbolt"
	"github.com/jeromedoucet/dahu/core/model"
)

func (i *inMemory) CreateJob(job *model.Job, ctx context.Context) (*model.Job, PersistenceError) {
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
		return nil, wrapError(err)
	}
}

func (i *inMemory) GetJob(id []byte, ctx context.Context) (*model.Job, PersistenceError) {
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
		return nil, wrapError(err)
	}
}

func (i *inMemory) GetJobs(ctx context.Context) ([]*model.Job, PersistenceError) {
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
		return nil, wrapError(err)
	}
}

func (i *inMemory) UpsertJobExecution(ctx context.Context, jobId string, execution *model.JobExecution) (*model.JobExecution, PersistenceError) {
	return nil, nil
}

func doFetchJobs(c *bolt.Cursor, jobs []*model.Job) ([]*model.Job, error) {
	res := jobs
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var job model.Job
		mErr := json.Unmarshal(v, &job)
		if mErr != nil {
			return nil, mErr
		} else {
			res = append(res, &job)
		}
	}
	return res, nil
}
