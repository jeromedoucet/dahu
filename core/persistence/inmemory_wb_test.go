package persistence

import (
	"context"
	"encoding/json"
	"testing"

	bolt "github.com/coreos/bbolt"
	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
)

// test enforcement of unicity of job
// test error on marshal

func TestCreateJobShouldReturnAnErrorWhenJobHasAnId(t *testing.T) {
	// given
	j := model.Job{Name: "test", Url: "github.com/test"}
	j.GenerateId()
	c := configuration.InitConf()

	ctx := context.Background()
	r := GetRepository(c)

	// when
	nj, err := r.CreateJob(&j, ctx)

	// then
	if nj != nil {
		t.Errorf(`expect to get no new job for a call on #CreateJob
		with a job that already have an id but got %+v`, nj)
	}
	if err == nil {
		t.Error(`expect to have an error when calling #CreateJob with
		a job that already have an id, but got nil`)
	}
	close(c.Close)
	r.WaitClose()
}

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

	// then
	if nj != nil {
		t.Errorf(`expect to get no new job for a call on #CreateJob
		when not bucket but got %+v`, nj)
	}
	if err == nil {
		t.Error(`expect to have an error when calling #CreateJob when
		no bucket, but got nil`)
	}
	close(c.Close)
	r.WaitClose()
}

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

	// then
	if err != nil {
		t.Errorf("expect to have no error when finding existing user, but got %s", err.Error())
	}
	if actualUser.Login != u.Login {
		t.Errorf("expect to get user %s but got %s", u.String(), actualUser.String())
	}
	close(c.Close)
	r.WaitClose()
}
