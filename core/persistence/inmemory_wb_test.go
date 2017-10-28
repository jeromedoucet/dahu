package persistence

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	bolt "github.com/coreos/bbolt"
	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
)

// todo test unicity for user

// test that we may not try to insert / create
// a job that already has an id
func TestCreateJobShouldReturnAnErrorWhenJobHasAnId(t *testing.T) {
	// given
	j := model.Job{Name: "test", Url: "github.com/test"}
	j.GenerateId()
	c := configuration.InitConf()

	ctx := context.Background()
	r := GetRepository(c)

	// when
	nj, err := r.CreateJob(&j, ctx)

	// close and remove the db
	close(c.Close)
	r.WaitClose()
	os.Remove(c.PersistenceConf.Name)

	// then
	if nj != nil {
		t.Errorf(`expect to get no new job for a call on #CreateJob
		with a job that already have an id but got %+v`, nj)
	}
	if err == nil {
		t.Error(`expect to have an error when calling #CreateJob with
		a job that already have an id, but got nil`)
	}

}

// test that the case when the bucket 'jobs' is missing is
// properly handle => an error is returned
func TestCreateJobShouldReturnAnErrorWhenNoBucket(t *testing.T) {
	// given
	j := model.Job{Name: "test", Url: "github.com/test"}
	c := configuration.InitConf()

	ctx := context.Background()
	rep := GetRepository(c)
	r, _ := rep.(*inMemory)
	r.db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("jobs"))
		return nil
	})

	// when
	nj, err := r.CreateJob(&j, ctx)

	// close and remove the db
	close(c.Close)
	r.WaitClose()
	os.Remove(c.PersistenceConf.Name)

	// then
	if nj != nil {
		t.Errorf(`expect to get no new job for a call on #CreateJob
		when not bucket but got %+v`, nj)
	}
	if err == nil {
		t.Error(`expect to have an error when calling #CreateJob when
		no bucket, but got nil`)
	}
}

// test the nominal case of #GetUser
func TestGetUserShouldReturnTheUserWhenItExists(t *testing.T) {
	// given
	u := model.User{Login: "test"}
	u.SetPassword([]byte("test_test_test_test"))
	c := configuration.InitConf()

	ctx := context.Background()
	rep := GetRepository(c)
	r, _ := rep.(*inMemory)
	// insertion of existing user
	r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		var data []byte
		data, _ = json.Marshal(u)
		b.Put([]byte(u.Login), data)
		return nil
	})

	// when
	actualUser, err := rep.GetUser([]byte(u.Login), ctx)

	// close and remove the db
	close(c.Close)
	r.WaitClose()
	os.Remove(c.PersistenceConf.Name)

	// then
	if err != nil {
		t.Errorf("expect to have no error when finding existing user, but got %s", err.Error())
	}
	if actualUser.Login != u.Login {
		t.Errorf("expect to get user %s but got %s", u.String(), actualUser.String())
	}
}

// test the case when the user is not found for #GetUser.
// => an error is returned
func TestGetUserShouldReturnAnErrorWhenItDoesntExist(t *testing.T) {
	// given
	u := model.User{Login: "test"}
	u.SetPassword([]byte("test_test_test_test"))
	c := configuration.InitConf()

	ctx := context.Background()
	rep := GetRepository(c)
	r, _ := rep.(*inMemory)

	// when
	actualUser, err := rep.GetUser([]byte(u.Login), ctx)

	// close and remove the db
	close(c.Close)
	r.WaitClose()
	os.Remove(c.PersistenceConf.Name)

	// then
	if err == nil {
		t.Error("expect to have an error when searching non-existing user, but got nil")
	}
	if actualUser != nil {
		t.Errorf("expect to get nil but got %s", actualUser.String())
	}
}

// test that the case where there is no bucket for user
// is properly handle at #GetUser
func TestGetUserShouldReturnAnErrorWhenNoBucket(t *testing.T) {
	// given
	u := model.User{Login: "test"}
	u.SetPassword([]byte("test_test_test_test"))
	c := configuration.InitConf()

	ctx := context.Background()
	rep := GetRepository(c)
	r, _ := rep.(*inMemory)
	r.db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("users"))
		return nil
	})

	// when
	actualUser, err := rep.GetUser([]byte(u.Login), ctx)

	// close and remove the db
	close(c.Close)
	r.WaitClose()
	os.Remove(c.PersistenceConf.Name)

	// then
	if err == nil {
		t.Error("expect to have an error when searching user without any users bucket, but got nil")
	}
	if actualUser != nil {
		t.Errorf("expect to get nil but got %s", actualUser.String())
	}
}
