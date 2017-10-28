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

// this mutex is used to ensure
// that inMemorySingleton is
// a singleton.
var oInMemory = &sync.Once{}

// used when stopping the db.
// it allow to wait for all transaction
// to be finished
var wgTransaction = &sync.WaitGroup{}

// this sync is used only when closing
// the current underlying db instance
var resetMutex = &sync.Mutex{}

var inMemorySingleton *inMemory

func getOrCreateInMemory(conf *configuration.Conf) Repository {
	oInMemory.Do(func() {
		createInMemory(conf)
	})
	return inMemorySingleton
}

func createInMemory(conf *configuration.Conf) {
	inMemorySingleton = new(inMemory)
	inMemorySingleton.conf = conf
	inMemorySingleton.rwMutex = &sync.Mutex{}
	inMemorySingleton.waitClose = &sync.WaitGroup{}
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
	// start lock for WaitClose func
	inMemorySingleton.waitClose.Add(1)
	go func() {
		// release lock for WaitClose func
		defer inMemorySingleton.waitClose.Done()
		<-inMemorySingleton.conf.Close
		// according to bbolt documentation
		// all transactions must be closed
		// before closing the db
		wgTransaction.Wait()
		inMemorySingleton.db.Close() // todo handle this error
		// reset the singletong
		oInMemory = &sync.Once{}
	}()
}

// In memory implementation
// of persistence.Repository.
// Use bbold as embedded db.
type inMemory struct {
	conf      *configuration.Conf
	db        *bolt.DB
	rwMutex   *sync.Mutex
	waitClose *sync.WaitGroup
}

func (i *inMemory) WaitClose() {
	i.waitClose.Wait()
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
		b.Put([]byte(job.Name), data)
		return nil
	})
	// if there is an error don't return the
	// job. This error must not be ignored
	if err == nil {
		return job, nil
	} else {
		return nil, err
	}

}

func (i *inMemory) GetJob(id []byte, ctx context.Context) (*model.Job, error) {
	panic("not implemented")
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
	// this increments will avoid
	// closing database while a transaction
	// is open
	wgTransaction.Add(1)
	defer wgTransaction.Done()
	select {
	case <-i.conf.Close:
		return errors.New("persistence >> the database is close or closing. Operation impossible.")
	default:
		// only one read/write transaction is allowed.
		i.rwMutex.Lock()
		defer i.rwMutex.Unlock()
		return i.db.Update(action)
	}

}

// doViewAction will ensure that the
// db is available and then execute the function
// inside a read transaction. An error is returned
// if an issue appears.
func (i *inMemory) doViewAction(action func(tx *bolt.Tx) error) error {
	// this increments will avoid
	// closing database while a transaction
	// is open
	wgTransaction.Add(1)
	defer wgTransaction.Done()
	select {
	case <-i.conf.Close:
		return errors.New("persistence >> the database is close or closing. Operation impossible.")
	default:
		return i.db.View(action)
	}

}
