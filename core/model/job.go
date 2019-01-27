package model

// this file source contains everything related
// to the jobs modelization. It means there is no
// code related to how jobs are executed here, but
// all data structure, from configuration to past
// execution.

import (
	"fmt"
	"strings"
	"time"
)

// configuration of a dahu job
type Job struct {
	Id              []byte         `json:"id"`              // id of the job
	Name            string         `json:"name"`            // simple label used for display
	GitConf         GitConfig      `json:"gitConfig"`       // repository configuration
	Steps           []Step         `json:"steps"`           // job steps execution
	Executions      []JobExecution `json:"executions"`      // list of past executions that are still available
	RemoveWorkspace bool           `json:"removeWorkspace"` // if true, the workspace is removed after every execution of the job
}

func (j *Job) GenerateId() error {
	id, err := generateId(j.Id)
	if err == nil {
		j.Id = id
	}
	return err
}

func (j *Job) IsValid() bool {
	if j.Name == "" || !j.GitConf.IsValid() {
		return false
	}
	return true
}

func (j *Job) String() string {
	return fmt.Sprintf("{Id:%s, Name:%s}", j.Id, j.Name)
}

func (j *Job) ToPublicModel() {
	j.GitConf.ToPublicModel()
}

// step of a Job. It is defined
// by an image and a command to run.
// Optionnally, some dependencies may
// be defined.
type Step struct {
	Name          string // display name of the step
	Image         Image  // the image that contains the needed dependencies for this step (node, golang, java, etc...)
	Envs          map[string]string
	Command       []string   // the command of the step
	MountingPoint string     // the place where the volume should be mounted. TODO think of a default value ?
	Services      []*Service // services that are needed for this step. For example a Database for an integration tests step.
}

// return Envs of the step and
// additional env entries bases on
// available services
func (s Step) ComputeEnvs() map[string]string {
	var res map[string]string
	if s.Envs == nil && len(s.Services) > 0 {
		res = make(map[string]string, len(s.Services))
	} else {
		res = s.Envs
	}
	for _, service := range s.Services {
		// for each services on step, we know that
		// it will be available through its name
		res[fmt.Sprintf("%s_HOST", service.Name)] = service.Name
	}
	return res
}

// Service that may needed for
// some step (for example integration tests).
// A name, the image and exposed port have to
// be defined
type Service struct {
	Name         string  // name under wich the service will be available during the step
	Image        Image   // container image
	ExposedPorts []*Port // exposed ports
}

// Container image. Contains
// its name and a registry ID, if
// needed (case of image not public on defaut registry)
type Image struct {
	Name       string // Name of the image
	RegistryId string // external key to a registry configuration
	Registry   *DockerRegistry
}

// ComputeName will return the Name of the image
// regarding of its registry. If there is none, the raw name
// is used, ${REGISTRY_URL}/${IMAGE_NAME} instead
func (i *Image) ComputeName() string {
	if i.Registry == nil {
		return i.Name
	} else {
		return fmt.Sprintf("%s/%s", i.Registry.Url, strings.TrimPrefix(strings.TrimSpace(i.Name), "/"))
	}
}

// Port exposed by a
// container. typically
// used by a StepDependency
// like a DataBase
type Port struct {
	Num       int
	Prototype string
}

// status for StepExecution
type ExecutionStatus string

const (
	Pending  ExecutionStatus = "pending"
	Running  ExecutionStatus = "running"
	Success  ExecutionStatus = "success"
	Failure  ExecutionStatus = "failure"
	Canceled ExecutionStatus = "canceled"
)

// contains everything related to
// one particular execution of a Job
type JobExecution struct {
	Id         string // the id of this execution Job. Used to update on particular execution
	BranchName string
	VolumeName string           // the name of the volume where the workspace is stored
	Steps      []*StepExecution // execution of step related to that job execution
	Date       time.Time        // the instant when the job execution has start
	Duration   time.Duration    // global duration of the job execution
}

func (j *JobExecution) GenerateId() error {
	id, err := generateId([]byte(j.Id))
	if err == nil {
		j.Id = string(id)
	}
	return err
}

// contains everything related to
// one execution of a step of a particular job execution
type StepExecution struct {
	Name     string
	Status   ExecutionStatus // status of the step execution
	Duration time.Duration   // global duration of the step execution
	Logs     string          // logs attached to the step
}

func (e StepExecution) IsSuccess() bool {
	return e.Status == Success
}
