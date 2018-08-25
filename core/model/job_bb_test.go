package model_test

import (
	"testing"

	"github.com/jeromedoucet/dahu/core/model"
)

func TestIsValidJobReturnTrue(t *testing.T) {
	// given
	httpAuth := model.HttpAuthConfig{Url: "http://some-domain/some-repo"}
	gitConf := model.GitConfig{HttpAuth: &httpAuth}
	j := model.Job{Name: "test", GitConf: gitConf}

	// when
	res := j.IsValid()

	// then
	if !res {
		t.Error("expect the job to be valid but it is not the case")
	}
}

func TestIsValidJobWithtInvalidGitConf(t *testing.T) {
	// given
	httpAuth := model.HttpAuthConfig{}
	gitConf := model.GitConfig{HttpAuth: &httpAuth}
	j := model.Job{Name: "test", GitConf: gitConf}

	// when
	res := j.IsValid()

	// then
	if res {
		t.Error("expect the job to be invalid but is valid")
	}
}

func TestJobIdGenerationShouldBeSuccessFullIfNoExistingId(t *testing.T) {
	// given
	j := new(model.Job)

	// when
	err := j.GenerateId()

	// then
	if err != nil {
		t.Errorf("Expect #GenerateId to return nil, but got %v", err)
	}
	if j.Id == nil {
		t.Errorf("expect the Id to have been generated, but is nil")
	}
}

func TestJobIdGenerationShouldBeInErrorIfExistingId(t *testing.T) {
	// given
	j := new(model.Job)
	j.Id = []byte("existingId")

	// when
	err := j.GenerateId()

	// then
	if err == nil {
		t.Errorf("Expect #GenerateId to return an error, but got nil")
	}
	if string(j.Id) != "existingId" {
		t.Errorf("expect the Id not to have changed, but got %s", string(j.Id))
	}
}

// case of StepExecution#IsSuccess that return true
func TestStepExecutionNominal(t *testing.T) {
	// given
	stepExecution := model.StepExecution{Status: model.Success}

	// when
	res := stepExecution.IsSuccess()

	// then
	if !res {
		t.Fatal("expect the stepExecution to be successful, but it is not")
	}
}

// case of StepExecution#IsSuccess that return false because of failure status
func TestStepExecutionFailed(t *testing.T) {
	// given
	stepExecution := model.StepExecution{Status: model.Failure}

	// when
	res := stepExecution.IsSuccess()

	// then
	if res {
		t.Fatal("expect the stepExecution to be failed, but it is not")
	}
}
