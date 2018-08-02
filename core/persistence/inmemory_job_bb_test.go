package persistence_test

import (
	"context"
	"testing"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
	"github.com/jeromedoucet/dahu/tests"
)

// test that we may not try to insert / create
// a job that already has an id
func TestCreateJobShouldReturnAnErrorWhenJobHasAnId(t *testing.T) {
	// given
	j := model.Job{Name: "test"}
	j.GenerateId()
	c := configuration.InitConf()

	ctx := context.Background()
	r := persistence.GetRepository(c)

	// when
	nj, err := r.CreateJob(&j, ctx)

	// close and remove the db
	tests.CleanPersistence(c)

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

// test the nominal case of #GetJob
func TestGetJobShouldReturnTheJobWhenItExists(t *testing.T) {
	// given
	j := model.Job{Name: "test"}
	j.GenerateId()
	c := configuration.InitConf()

	ctx := context.Background()
	tests.InsertObject(c, []byte("jobs"), []byte(j.Id), j)
	rep := persistence.GetRepository(c)

	// when
	actualJob, err := rep.GetJob([]byte(j.Id), ctx)

	// close and remove the db
	tests.CleanPersistence(c)

	// then
	if err != nil {
		t.Errorf("expect to have no error when finding existing job, but got %s", err.Error())
	}
	if string(actualJob.Id) != string(j.Id) || actualJob.Name != j.Name {
		t.Errorf("expect to get user %s but got %s", j.String(), actualJob.String())
	}
}

// test the nominal case of #GetJob
func TestGetJobShouldReturnAnErrorWhenItDoesntExists(t *testing.T) {
	// given
	j := model.Job{Name: "test"}
	j.GenerateId()
	c := configuration.InitConf()

	ctx := context.Background()
	tests.InsertObject(c, []byte("jobs"), []byte(j.Id), j)
	rep := persistence.GetRepository(c)

	// when
	actualJob, err := rep.GetJob([]byte("test2"), ctx)

	// close and remove the db
	tests.CleanPersistence(c)

	// then
	if err == nil {
		t.Error("expect to have an error when searching non-existing job, but got nil")
	}
	if actualJob != nil {
		t.Errorf("expect to get nil but got %s", actualJob.String())
	}
}
