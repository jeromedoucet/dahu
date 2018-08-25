package model

// this file source contains everything related
// to the jobs modelization. It means there is no
// code related to how jobs are executed here, but
// all data structure, from configuration to past
// execution.

import (
	"fmt"
	"time"
)

// configuration of a dahu job
type Job struct {
	Id         []byte         `json:"id"`         // id of the job
	Name       string         `json:"name"`       // simple label used for display
	GitConf    GitConfig      `json:"gitConfig"`  // repository configuration
	Steps      []Step         `json:"steps"`      // job steps execution
	Executions []JobExecution `json:"executions"` // list of past executions that are still available
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
	Name          string     // display name of the step
	Image         Image      // the image that contains the needed dependencies for this step (node, golang, java, etc...)
	Command       string     // the command of the step
	MountingPoint string     // the place where the volume should be mounted. TODO think of a default value ?
	Services      []Services // services that are needed for this step. For example a Database for an integration tests step.
}

// Services that may needed for
// some step (for example integration tests).
// A name, the image and exposed port have to
// be defined
type Services struct {
	Name         string // name under wich the service will be available during the step
	Image        Image  // container image
	ExposedPorts []Port // exposed ports
}

// Container image. Contains
// its name and a registry ID, if
// needed (case of image not public on defaut registry)
type Image struct {
	Name       string // Name of the image
	RegistryId string // external key to a registry configuration
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
	Pending ExecutionStatus = "pending"
	Running ExecutionStatus = "running"
	Success ExecutionStatus = "success"
	Failure ExecutionStatus = "failure"
	Aborted ExecutionStatus = "aborted"
)

// contains everything related to
// one particular execution of a Job
type JobExecution struct {
	Id         string          // the id of this execution Job. Used to update on particular execution
	VolumeName string          // the name of the volume where the workspace is stored
	Steps      []StepExecution // execution of step related to that job execution
	Date       time.Time       // the instant when the job execution has start
	Duration   time.Duration   // global duration of the job execution
	// todo trigger ? User Name ?
}

// contains everything related to
// one execution of a step of a particular job execution
type StepExecution struct {
	Name          string
	Status        ExecutionStatus // status of the step execution
	Duration      time.Duration   // global duration of the step execution
	LogVolumeName string          // the name where the logs are stored
}

func (e StepExecution) IsSuccess() bool {
	return e.Status == Success
}
