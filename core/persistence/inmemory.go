package persistence

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

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
	inMemorySingleton.close = make(chan interface{})
	inMemorySingleton.conf = conf
	inMemorySingleton.rwMutex = &sync.RWMutex{}
	newDb := isNewDb(conf.PersistenceConf.Name)
	db, _ := bolt.Open(conf.PersistenceConf.Name, 0600, nil) // todo handle this error
	dbInitialization(db)
	if newDb {
		defaultUserCreation(db)
	}
	inMemorySingleton.db = db
	go waitForGracefullShutdown()
}

func isNewDb(dbName string) bool {
	ex, _ := os.Executable()
	dir := path.Dir(ex)
	_, err := os.Stat(dir + "/" + dbName)
	return os.IsNotExist(err)
}

func dbInitialization(db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		return createBucketsIfNeeded(tx)
	})
}

// will insert the default user.
// should only be used if a previous call
// of isNewDb has returned true !
func defaultUserCreation(db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		u := model.User{Login: "dahu"}
		u.SetPassword([]byte("dahuDefaultPassword"))
		b, _ := tx.CreateBucketIfNotExists([]byte("users"))
		var data []byte
		data, _ = json.Marshal(u)
		b.Put([]byte(u.Login), data)
		return nil
	})
}

func waitForGracefullShutdown() {
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
	close(inMemorySingleton.close)
	inMemorySingleton = nil
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
	close   chan interface{}
}

func (i *inMemory) WaitClose() {
	<-i.close
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
