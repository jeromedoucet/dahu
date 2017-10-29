package tests

import (
	"encoding/json"
	"fmt"
	"os"

	bolt "github.com/coreos/bbolt"
	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/persistence"
)

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
