package model_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/tests"
)

// test the behavior of the run module when
// executing a Job with a failure result.
// We must check that the correct status
// is return and all logs streamed
func TestRunStartWithFailure(t *testing.T) {
	// given
	buf := new(bytes.Buffer)
	params := model.RunParams{
		Id:    "test-1",
		Image: "dahuci/job-test",
		Env: map[string]string{
			"REPO_URL": "git@github.com:jeromedoucet/dahu-images.git",
			"STATUS":   "failure",
		},
		OutputWriter: buf,
	}
	defer tests.RemoveContainer(params.ContainerName())
	r := model.NewRun(params)

	// when
	r.Start(context.Background())

	<-r.Done()
	// then
	if r.Status() != model.FAILURE {
		t.Errorf("Expect FAILURE state (%d), got : %d", model.FAILURE, r.Status())
	}
	if buf.String() != "Failure\n" {
		t.Errorf("Expect 'Failure' in output writer, got : %#v", buf.String())
	}
}

// test the cancelation after a given time out
func TestRunStartWithTimeOut(t *testing.T) {
	// given
	buf := new(bytes.Buffer)
	params := model.RunParams{
		Id:    "test-2",
		Image: "dahuci/job-test",
		Env: map[string]string{
			"REPO_URL": "git@github.com:jeromedoucet/dahu-images.git",
			"STATUS":   "timeout",
		},
		OutputWriter: buf,
		TimeOut:      time.Second * 1,
	}
	defer tests.RemoveContainer(params.ContainerName())
	r := model.NewRun(params)

	// when
	r.Start(context.Background())

	<-r.Done()
	// then
	if r.Status() != model.FAILURE {
		t.Errorf("Expect FAILURE state (%d), got : %d", model.FAILURE, r.Status())
	}
	if buf.String() != "Time out" {
		t.Errorf("Expect 'Time out' in output writer, got : %#v", buf.String())
	}
}

// test the behavior of the run module when
// executing a Job with a success result.
func TestRunStartWithSuccess(t *testing.T) {
	// given
	buf := new(bytes.Buffer)
	params := model.RunParams{
		Id:    "test-3",
		Image: "dahuci/job-test",
		Env: map[string]string{
			"REPO_URL": "git@github.com:jeromedoucet/dahu-images.git",
			"STATUS":   "success",
		},
		OutputWriter: buf,
	}
	defer tests.RemoveContainer(params.ContainerName())
	r := model.NewRun(params)

	// when
	r.Start(context.Background())

	<-r.Done()
	// then
	if r.Status() != model.SUCCESS {
		t.Errorf("Expect SUCCESS state (%d), got : %d", model.SUCCESS, r.Status())
	}
	if buf.String() != "Success\n" {
		t.Errorf("Expect 'Success' in output writer, got : %#v", buf.String())
	}
}

// test the behavior of the run module when
// starting a Run twice.
func TestRunStartTwiceShouldReturnError(t *testing.T) {
	// given
	buf := new(bytes.Buffer)
	params := model.RunParams{
		Id:    "test-4",
		Image: "dahuci/job-test",
		Env: map[string]string{
			"REPO_URL": "git@github.com:jeromedoucet/dahu-images.git",
			"STATUS":   "success",
		},
		OutputWriter: buf,
	}
	defer tests.RemoveContainer(params.ContainerName())
	r := model.NewRun(params)

	// when
	err1 := r.Start(context.Background())
	err2 := r.Start(context.Background())

	<-r.Done()
	// then
	if err1 != nil {
		t.Errorf("Expect the first call to return no error, but got : %v", err1.Error())
	}
	if err2 == nil {
		t.Errorf("Expect the second call to return an error, but got : nil")
	}
	if r.Status() != model.SUCCESS {
		t.Errorf("Expect SUCCESS state (%d), got : %d", model.SUCCESS, r.Status())
	}
	if buf.String() != "Success\n" {
		t.Errorf("Expect 'Success' in output writer, got : %#v", buf.String())
	}
}

// test a cancelation of a run
func TestRunStartWithCancelation(t *testing.T) {
	// given
	buf := new(bytes.Buffer)
	params := model.RunParams{
		Id:    "test-5",
		Image: "dahuci/job-test",
		Env: map[string]string{
			"REPO_URL": "git@github.com:jeromedoucet/dahu-images.git",
			"STATUS":   "timeout",
		},
		OutputWriter: buf,
	}
	defer tests.RemoveContainer(params.ContainerName())
	r := model.NewRun(params)

	// when
	err := r.Start(context.Background())
	r.Cancel()

	<-r.Done()
	// then
	if err != nil {
		t.Errorf("Expect the first call to return no error, but got : %v", err.Error())
	}
	if r.Status() != model.CANCELED {
		t.Errorf("Expect CANCELED state (%d), got : %d", model.CANCELED, r.Status())
	}
}

// test that we can not cancel a run
// before starting it
func TestRunCancelationShouldFailWhenNotStarted(t *testing.T) {
	// given
	buf := new(bytes.Buffer)
	params := model.RunParams{
		Id:    "test-6",
		Image: "dahuci/job-test",
		Env: map[string]string{
			"REPO_URL": "git@github.com:jeromedoucet/dahu-images.git",
			"STATUS":   "timeout",
		},
		OutputWriter: buf,
	}
	r := model.NewRun(params)

	// when
	err := r.Cancel()

	// then
	if err == nil {
		t.Error("Expect the cancel call to return an error, but got nil")
	}
}

// test that we can not cancel a run after it has finished
func TestRunCancelationShouldFailWhenFinished(t *testing.T) {
	// given
	buf := new(bytes.Buffer)
	params := model.RunParams{
		Id:    "test-7",
		Image: "dahuci/job-test",
		Env: map[string]string{
			"REPO_URL": "git@github.com:jeromedoucet/dahu-images.git",
			"STATUS":   "success",
		},
		OutputWriter: buf,
	}

	defer tests.RemoveContainer(params.ContainerName())
	r := model.NewRun(params)

	// when
	r.Start(context.Background())
	<-r.Done()
	err := r.Cancel()

	// then
	if err == nil {
		t.Error("Expect the cancel call to return an error, but got nil")
	}
}

// todo test with existing container name