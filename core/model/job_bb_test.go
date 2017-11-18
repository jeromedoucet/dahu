package model_test

import (
	"os/exec"
	"testing"

	"github.com/jeromedoucet/dahu/core/model"
)

func TestIsValidJobRunReturnTrue(t *testing.T) {
	// given
	jr := model.JobRun{Status: model.CREATED, ContainerName: "test"}

	// when
	res := jr.IsValid()

	// then
	if !res {
		t.Error("expect a RunStatus without issue to be valid but is invalid")
	}
}

func TestIsValidJobRunWithoutContainerName(t *testing.T) {
	// given
	jr := model.JobRun{Status: model.CREATED}

	// when
	res := jr.IsValid()

	// then
	if res {
		t.Error("expect a RunStatus without ContainerName to be invalid but is valid")
	}
}

func TestIsValidJobRunWithoutStatus(t *testing.T) {
	// given
	jr := model.JobRun{ContainerName: "test"}

	// when
	res := jr.IsValid()

	// then
	if res {
		t.Error("expect a RunStatus without status to be invalid but is valid")
	}
}

// This is the test of the roll behavior of
// jobRun inside the Job. The slice must not
// grow and when appending a new JobRun,
// the older are push to the end. Every time
// it happen, the oldest is removed and the related container
// is removed
func TestAppendJobRunShouldRollJobRun(t *testing.T) {
	// given

	// run a docker container that should be removed
	cmd := exec.Command("docker", "run", "--name", "dahu-test", "hello-world")
	cmd.Run()

	j := model.Job{JobRuns: []*model.JobRun{
		&model.JobRun{Id: []byte("5")},
		&model.JobRun{Id: []byte("4")},
		&model.JobRun{Id: []byte("3")},
		&model.JobRun{Id: []byte("2")},
		&model.JobRun{Id: []byte("1"), ContainerName: "dahu-test"},
	}}
	jr := model.JobRun{Id: []byte("6")}

	// when
	j.AppendJobRun(&jr)

	// then
	if string(j.JobRuns[0].Id) != string(jr.Id) {
		t.Errorf("expect the first JobRun to be %s", string(jr.Id))
	}
	if string(j.JobRuns[1].Id) != "5" {
		t.Error("expect the second JobRun to be 5")
	}
	if string(j.JobRuns[2].Id) != "4" {
		t.Error("expect the third JobRun to be 4")
	}
	if string(j.JobRuns[3].Id) != "3" {
		t.Error("expect the fourth JobRun to be 3")
	}
	if string(j.JobRuns[4].Id) != "2" {
		t.Error("expect the last JobRun to be 2")
	}

	// check the docker container is removed
	cmd = exec.Command("docker", "rm", "-f", "dahu-test")
	err := cmd.Run()
	if err == nil {
		t.Error("expect the container to be removed but it was not")
	}
}

// this test will ensure that #AppendJobRun behavior
// is correct even when the inner slice is not
// initiate
func TestAppendJobRunShouldCreateTheSliceWhenNil(t *testing.T) {
	// given
	j := model.Job{}
	jr := model.JobRun{Id: []byte("6")}

	// when
	j.AppendJobRun(&jr)

	// then
	if string(j.JobRuns[0].Id) != string(jr.Id) {
		t.Errorf("expect the first JobRun to be %s", string(jr.Id))
	}
	if j.JobRuns[1] != nil {
		t.Error("expect the second JobRun to be nil")
	}
	if j.JobRuns[2] != nil {
		t.Error("expect the third JobRun to be nil")
	}
	if j.JobRuns[3] != nil {
		t.Error("expect the fourth JobRun to be nil")
	}
	if j.JobRuns[4] != nil {
		t.Error("expect the last JobRun to be nil")
	}
}

// todo test when jobRun not valid
