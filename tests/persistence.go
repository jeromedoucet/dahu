package tests

import (
	"encoding/json"
	"fmt"
	"os"

	bolt "github.com/coreos/bbolt"
	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/persistence"
)

type MockRepository struct {
	persistence.Repository
	CreateJobRunCount int
	UpdateJobRunCount int
}

// will Close persistence layer (if needed)
// and then remove all data within.
// should be used at the end of a test
func CleanPersistence(conf *configuration.Conf) {
	ClosePersistence(conf)
	DeletePersistence(conf)
}

// Will Close persistence layer without deleting
// data. Is usefull when it is required to make
// some check at the end of a test with a persistence
// system that don't allow concurrent access (bbolt)
//
// If this function is used, a call to #DeletePersistence
// may be required to clean everything before the next test
func ClosePersistence(conf *configuration.Conf) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in CleanPersistence", r)
		}
	}()
	close(conf.Close)
	rep := persistence.GetRepository(conf)
	rep.WaitClose()
}

func DeletePersistence(conf *configuration.Conf) {
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

func ObjectExist(conf *configuration.Conf, bucketName, key []byte) bool {
	res := false
	db, _ := bolt.Open(conf.PersistenceConf.Name, 0600, nil)
	db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketName).Get(key)
		if data != nil {
			res = true
		}
		return nil
	})
	db.Close()
	return res
}

func DeleteBucket(conf *configuration.Conf, bucketName []byte) {
	db, _ := bolt.Open(conf.PersistenceConf.Name, 0600, nil)
	db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(bucketName)
		return nil
	})
	db.Close()
}
