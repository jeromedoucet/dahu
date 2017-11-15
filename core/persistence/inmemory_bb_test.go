package persistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
	"github.com/jeromedoucet/dahu/tests"
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

// test that the case when the bucket 'jobs' is missing is
// properly handle => an error is returned
func TestCreateJobShouldReturnAnErrorWhenNoBucket(t *testing.T) {
	// todo move it in a 'white box' to make it possible (with no deadlock)
	t.SkipNow()
	// given
	j := model.Job{Name: "test", Url: "github.com/test"}
	c := configuration.InitConf()

	ctx := context.Background()
	rep := persistence.GetRepository(c)
	tests.DeleteBucket(c, []byte("jobs"))

	// when
	nj, err := rep.CreateJob(&j, ctx)

	// close and remove the db
	tests.CleanPersistence(c)

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

// test the nominal case of #GetUser
func TestGetUserShouldReturnTheUserWhenItExists(t *testing.T) {
	// given
	u := model.User{Login: "test"}
	u.SetPassword([]byte("test_test_test_test"))
	c := configuration.InitConf()

	ctx := context.Background()
	tests.InsertObject(c, []byte("users"), []byte(u.Login), u)
	rep := persistence.GetRepository(c)

	// when
	actualUser, err := rep.GetUser([]byte(u.Login), ctx)

	// close and remove the db
	tests.CleanPersistence(c)

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
	rep := persistence.GetRepository(c)

	// when
	actualUser, err := rep.GetUser([]byte(u.Login), ctx)

	// close and remove the db
	tests.CleanPersistence(c)

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
	// todo move it in a 'white box' to make it possible (with no deadlock)
	t.SkipNow()
	// given
	u := model.User{Login: "test"}
	u.SetPassword([]byte("test_test_test_test"))
	c := configuration.InitConf()
	tests.DeleteBucket(c, []byte("users"))

	ctx := context.Background()
	rep := persistence.GetRepository(c)

	// when
	actualUser, err := rep.GetUser([]byte(u.Login), ctx)

	// close and remove the db
	tests.CleanPersistence(c)

	// then
	if err == nil {
		t.Error("expect to have an error when searching user without any users bucket, but got nil")
	}
	if actualUser != nil {
		t.Errorf("expect to get nil but got %s", actualUser.String())
	}
}

// test the nominal case of #CreateJobRun
func TestCreateJobRunShouldUpdateAJob(t *testing.T) {
	// given
	j := model.Job{Name: "test"}
	j.GenerateId()
	jr := model.JobRun{ContainerName: "test", Status: model.RUNNING, StartTime: time.Now()}
	c := configuration.InitConf()

	ctx := context.Background()
	tests.InsertObject(c, []byte("jobs"), []byte(j.Id), j)
	rep := persistence.GetRepository(c)

	// when
	actualJobRun, err := rep.CreateJobRun(&jr, []byte(j.Id), ctx)

	job, _ := rep.GetJob([]byte(j.Id), ctx)
	// close and remove the db
	tests.CleanPersistence(c)

	// then
	if err != nil {
		t.Errorf("expect to have no error when creating a jobRun on an existing Job, but got %s", err.Error())
	}
	if string(actualJobRun.Id) == "" {
		t.Error("expect the new jobRun to have an id, but is empty.")
	}
	if string(job.JobRuns[0].Id) != string(actualJobRun.Id) {
		t.Errorf("expect the jobRun to be the first on the Job but got %+v", job.JobRuns[0])
	}
}

// test #CreateJobRun when JobRun is invalid
func TestCreateJobRunShouldFailIfJobRunInvalid(t *testing.T) {
	// given
	j := model.Job{Name: "test"}
	j.GenerateId()
	jr := model.JobRun{Status: model.RUNNING, StartTime: time.Now()}
	c := configuration.InitConf()

	ctx := context.Background()
	tests.InsertObject(c, []byte("jobs"), []byte(j.Id), j)
	rep := persistence.GetRepository(c)

	// when
	actualJobRun, err := rep.CreateJobRun(&jr, []byte(j.Id), ctx)

	job, _ := rep.GetJob([]byte(j.Id), ctx)
	// close and remove the db
	tests.CleanPersistence(c)

	// then
	if err == nil {
		t.Fatal("expect to have error when trying to create an invalid jobRun , but got nil")
	}
	if actualJobRun != nil {
		t.Fatalf("expect the new jobRun to be nil, but go %+v", actualJobRun)
	}
	if len(job.JobRuns) != 0 {
		t.Errorf("expect the jobRuns to be empty but got %+v", job.JobRuns[0])
	}
}

// test #CreateJobRun when JobRun is invalid
func TestCreateJobRunShouldFailIfNoJob(t *testing.T) {
	// given
	j := model.Job{Name: "test"}
	j.GenerateId()
	jr := model.JobRun{Status: model.RUNNING, StartTime: time.Now()}
	c := configuration.InitConf()

	ctx := context.Background()
	rep := persistence.GetRepository(c)

	// when
	actualJobRun, err := rep.CreateJobRun(&jr, []byte(j.Id), ctx)

	// close and remove the db
	tests.CleanPersistence(c)

	// then
	if err == nil {
		t.Fatal("expect to have error when trying to create a jobRun on a non existing Job, but got nil")
	}
	if actualJobRun != nil {
		t.Fatalf("expect the new jobRun to be nil, but go %+v", actualJobRun)
	}
}
// todo test id already exist
