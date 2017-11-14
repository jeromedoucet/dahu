package model_test

import (
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

func TestAppendJobRunShouldRollJobRun(t *testing.T) {
	// given
	j := model.Job{JobRuns: []*model.JobRun{
		&model.JobRun{Id: []byte("5")},
		&model.JobRun{Id: []byte("4")},
		&model.JobRun{Id: []byte("3")},
		&model.JobRun{Id: []byte("2")},
		&model.JobRun{Id: []byte("1")},
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
}

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
