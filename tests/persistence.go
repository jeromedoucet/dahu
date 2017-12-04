package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	bolt "github.com/coreos/bbolt"
	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
)

type MockRepository struct {
	persistence.Repository
	CreateJobRunCount int
	UpdateJobRunCount int
	CreateJobRuns     []model.JobRun
	UpdateJobRuns     []model.JobRun
}

func (m *MockRepository) CreateJobRun(jobRun *model.JobRun, jobId []byte, ctx context.Context) (*model.JobRun, error) {
	if m.CreateJobRuns == nil {
		m.CreateJobRuns = make([]model.JobRun, 0)
	}
	m.CreateJobRuns = append(m.CreateJobRuns, *jobRun)
	m.CreateJobRunCount++
	return jobRun, nil
}

func (m *MockRepository) UpdateJobRun(jobRun *model.JobRun, jobId []byte, ctx context.Context) (*model.JobRun, error) {
	if m.UpdateJobRuns == nil {
		m.UpdateJobRuns = make([]model.JobRun, 0)
	}
	m.UpdateJobRuns = append(m.UpdateJobRuns, *jobRun)
	m.UpdateJobRunCount++
	return jobRun, nil
}

// this package is a collection
// of functions used in this project tests.
// for the persistence layer.

// will clean the data inside persistence
// must be done after a test has run.
func CleanPersistence(conf *configuration.Conf) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in CleanPersistence", r)
		}
	}()
	close(conf.Close)
	rep := persistence.GetRepository(conf)
	rep.WaitClose()
	os.Remove(conf.PersistenceConf.Name)
}

// insert some objet inside a given bucket. Create the bucket
// if needed
func InsertObject(conf *configuration.Conf, bucketName, key []byte, object interface{}) {
	db, _ := bolt.Open(conf.PersistenceConf.Name, 0600, nil)
	db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists(bucketName)
		var data []byte
		data, _ = json.Marshal(object)
		b.Put(key, data)
		return nil
	})
	db.Close()
}

func DeleteBucket(conf *configuration.Conf, bucketName []byte) {
	db, _ := bolt.Open(conf.PersistenceConf.Name, 0600, nil)
	db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(bucketName)
		return nil
	})
	db.Close()
}
