package run_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
	"github.com/jeromedoucet/dahu/core/run"
	"github.com/jeromedoucet/dahu/tests"
)

type mockRepository struct {
	persistence.Repository
	createJobRunCount int
	updateJobRunCount int
	createJobRuns     []model.JobRun
	updateJobRuns     []model.JobRun
}

func (m *mockRepository) CreateJobRun(jobRun *model.JobRun, jobId []byte, ctx context.Context) (*model.JobRun, error) {
	if m.createJobRuns == nil {
		m.createJobRuns = make([]model.JobRun, 0)
	}
	m.createJobRuns = append(m.createJobRuns, *jobRun)
	m.createJobRunCount++
	return jobRun, nil
}

func (m *mockRepository) UpdateJobRun(jobRun *model.JobRun, jobId []byte, ctx context.Context) (*model.JobRun, error) {
	if m.updateJobRuns == nil {
		m.updateJobRuns = make([]model.JobRun, 0)
	}
	m.updateJobRuns = append(m.updateJobRuns, *jobRun)
	m.updateJobRunCount++
	return jobRun, nil
}

// test the behavior of the run module when
// executing a Job with a failure result.
// We must check that the correct status
// is return and all logs streamed
func TestProcessStartWithFailure(t *testing.T) {
	// given
	buf := new(bytes.Buffer)
	params := run.ProcessParams{
		Id:    "test-1",
		Image: "dahuci/job-test",
		Env: map[string]string{
			"REPO_URL": "git@github.com:jeromedoucet/dahu-images.git",
			"STATUS":   "failure",
		},
		OutputWriter: buf,
	}
	defer tests.RemoveContainer(params.ContainerName())
	m := new(mockRepository)
	r := run.NewProcess(params, m)

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
	if m.createJobRunCount != 1 {
		t.Error("expect #CreateJobRun to have been called one time, but it has not been called")
	}
	if m.createJobRuns[0].Status != model.RUNNING {
		t.Errorf("expect the jobRun status when created to be %d, but got %d", model.RUNNING, m.createJobRuns[0].Status)
	}
	if m.updateJobRunCount != 1 {
		t.Error("expect #UpdateJobRun to have been called one time, but it has not been called")
	}
	if m.updateJobRuns[0].Status != model.FAILURE {
		t.Errorf("expect the jobRun status when created to be %d, but got %d", model.FAILURE, m.updateJobRuns[0].Status)
	}
}

// test the cancelation after a given time out
func TestProcessStartWithTimeOut(t *testing.T) {
	// given
	buf := new(bytes.Buffer)
	params := run.ProcessParams{
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
	m := new(mockRepository)
	r := run.NewProcess(params, m)

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
	if m.createJobRunCount != 1 {
		t.Error("expect #CreateJobRun to have been called one time, but it has not been called")
	}
	if m.createJobRuns[0].Status != model.RUNNING {
		t.Errorf("expect the jobRun status when created to be %d, but got %d", model.RUNNING, m.createJobRuns[0].Status)
	}
	if m.updateJobRunCount != 1 {
		t.Error("expect #UpdateJobRun to have been called one time, but it has not been called")
	}
	if m.updateJobRuns[0].Status != model.FAILURE {
		t.Errorf("expect the jobRun status when created to be %d, but got %d", model.FAILURE, m.updateJobRuns[0].Status)
	}
}

// test the behavior of the run module when
// executing a Job with a success result.
func TestProcessStartWithSuccess(t *testing.T) {
	// given
	buf := new(bytes.Buffer)
	params := run.ProcessParams{
		Id:    "test-3",
		Image: "dahuci/job-test",
		Env: map[string]string{
			"REPO_URL": "git@github.com:jeromedoucet/dahu-images.git",
			"STATUS":   "success",
		},
		OutputWriter: buf,
	}
	defer tests.RemoveContainer(params.ContainerName())
	m := new(mockRepository)
	r := run.NewProcess(params, m)

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
	if m.createJobRunCount != 1 {
		t.Error("expect #CreateJobRun to have been called one time, but it has not been called")
	}
	if m.createJobRuns[0].Status != model.RUNNING {
		t.Errorf("expect the jobRun status when created to be %d, but got %d", model.RUNNING, m.createJobRuns[0].Status)
	}
	if m.updateJobRunCount != 1 {
		t.Error("expect #UpdateJobRun to have been called one time, but it has not been called")
	}
	if m.updateJobRuns[0].Status != model.SUCCESS {
		t.Errorf("expect the jobRun status when created to be %d, but got %d", model.SUCCESS, m.updateJobRuns[0].Status)
	}
}

// test the behavior of the run module when
// starting a Run twice.
func TestProcessStartTwiceShouldReturnError(t *testing.T) {
	// given
	buf := new(bytes.Buffer)
	params := run.ProcessParams{
		Id:    "test-4",
		Image: "dahuci/job-test",
		Env: map[string]string{
			"REPO_URL": "git@github.com:jeromedoucet/dahu-images.git",
			"STATUS":   "success",
		},
		OutputWriter: buf,
	}
	defer tests.RemoveContainer(params.ContainerName())
	r := run.NewProcess(params, new(mockRepository))

	// when
	_, err1 := r.Start(context.Background())
	_, err2 := r.Start(context.Background())

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
func TestProcessStartWithCancelation(t *testing.T) {
	// given
	buf := new(bytes.Buffer)
	params := run.ProcessParams{
		Id:    "test-5",
		Image: "dahuci/job-test",
		Env: map[string]string{
			"REPO_URL": "git@github.com:jeromedoucet/dahu-images.git",
			"STATUS":   "timeout",
		},
		OutputWriter: buf,
	}
	defer tests.RemoveContainer(params.ContainerName())
	m := new(mockRepository)
	r := run.NewProcess(params, m)

	// when
	_, err := r.Start(context.Background())
	r.Cancel(context.Background())

	<-r.Done()
	// then
	if err != nil {
		t.Errorf("Expect the first call to return no error, but got : %v", err.Error())
	}
	if r.Status() != model.CANCELED {
		t.Errorf("Expect CANCELED state (%d), got : %d", model.CANCELED, r.Status())
	}
	if m.createJobRunCount != 1 {
		t.Error("expect #CreateJobRun to have been called one time, but it has not been called")
	}
	if m.createJobRuns[0].Status != model.RUNNING {
		t.Errorf("expect the jobRun status when created to be %d, but got %d", model.RUNNING, m.createJobRuns[0].Status)
	}
	if m.updateJobRunCount != 1 {
		t.Error("expect #UpdateJobRun to have been called one time, but it has not been called")
	}
	if m.updateJobRuns[0].Status != model.CANCELED {
		t.Errorf("expect the jobRun status when created to be %d, but got %d", model.CANCELED, m.updateJobRuns[0].Status)
	}
}

// test that we can not cancel a run
// before starting it
func TestProcessCancelationShouldFailWhenNotStarted(t *testing.T) {
	// given
	buf := new(bytes.Buffer)
	params := run.ProcessParams{
		Id:    "test-6",
		Image: "dahuci/job-test",
		Env: map[string]string{
			"REPO_URL": "git@github.com:jeromedoucet/dahu-images.git",
			"STATUS":   "timeout",
		},
		OutputWriter: buf,
	}
	r := run.NewProcess(params, new(mockRepository))

	// when
	err := r.Cancel(context.Background())

	// then
	if err == nil {
		t.Error("Expect the cancel call to return an error, but got nil")
	}
}

// test that we can not cancel a run after it has finished
func TestProcessCancelationShouldFailWhenFinished(t *testing.T) {
	// given
	buf := new(bytes.Buffer)
	params := run.ProcessParams{
		Id:    "test-7",
		Image: "dahuci/job-test",
		Env: map[string]string{
			"REPO_URL": "git@github.com:jeromedoucet/dahu-images.git",
			"STATUS":   "success",
		},
		OutputWriter: buf,
	}

	defer tests.RemoveContainer(params.ContainerName())
	r := run.NewProcess(params, new(mockRepository))

	// when
	r.Start(context.Background())
	<-r.Done()
	err := r.Cancel(context.Background())

	// then
	if err == nil {
		t.Error("Expect the cancel call to return an error, but got nil")
	}
}

// todo test with existing container name
// todo test CreateJobRun is called with error
// todo test correct handle on CreateJobRun and UpdateJobRun
