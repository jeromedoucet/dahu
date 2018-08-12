package model

import (
	"fmt"
)

type Job struct {
	Id      []byte    `json:"id"`
	Name    string    `json:"name"`
	GitConf GitConfig `json:"gitConfig"`
	Steps   []Step    `json:"steps"`
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
