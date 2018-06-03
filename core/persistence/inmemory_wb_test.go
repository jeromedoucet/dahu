package persistence

import (
	"errors"
	"testing"

	bolt "github.com/coreos/bbolt"
)

type txMock struct {
	bucketToReject string
}

func (tx *txMock) CreateBucketIfNotExists(name []byte) (b *bolt.Bucket, err error) {
	if string(name) == tx.bucketToReject {
		err = errors.New("some error")
	}
	return
}

var dbName string = "inMemoryWb"

/*
* Test case of createBucketIfNeeded when the users bucket creation
* fail. Should return an error
 */
func TestCreateBucketsIfNeededUsersBucketError(t *testing.T) {
	// given
	tx := new(txMock)
	tx.bucketToReject = "users"
	expectedErrorMsg := "ERROR >> user bucket creation failed : some error"

	// when
	err := createBucketsIfNeeded(tx)

	// then
	if err == nil {
		t.Error("expect to have an error, but got nil")
	}

	if err.Error() != expectedErrorMsg {
		t.Errorf("expect to have an error with text %s, but got %s", expectedErrorMsg, err.Error())
	}
}

/*
* Test case of createBucketIfNeeded when the jobs bucket creation
* fail. Should return an error
 */
func TestCreateBucketsIfNeededJObsBucketError(t *testing.T) {
	// given
	tx := new(txMock)
	tx.bucketToReject = "jobs"
	expectedErrorMsg := "ERROR >> job bucket creation failed : some error"

	// when
	err := createBucketsIfNeeded(tx)

	// then
	if err == nil {
		t.Error("expect to have an error, but got nil")
	}

	if err.Error() != expectedErrorMsg {
		t.Errorf("expect to have an error with text %s, but got %s", expectedErrorMsg, err.Error())
	}

}
