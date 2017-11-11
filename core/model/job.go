package model

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// use by request that
// will start a job run.
type RunRequest struct {
	OpenWs bool `json:"openWs"`
}

// A job run store informations on a run
// for a given Job
type JobRun struct {
	Id            []byte    `json:"id"`
	ContainerName string    `json:"containerName"`
	Status        RunStatus `json: "runStatus"`
	StartTime     time.Time `json: "startTime"`
	EndTime       time.Time `json:"endTime"`
}

// generate an Id for the JobRun.
// if it already exist, return an error
func (j *JobRun) GenerateId() error {
	id, err := generateId(j.Id)
	if err == nil {
		j.Id = id
	}
	return err
}

// Configuration detail of
// a Job
type JobConfiguration struct {
	NbRunBackup int `json:"nbRunBackup"` // the number of Run result that are kept
}

type Job struct {
	Id        []byte            `json:"id"` // todo change it into a []byte
	Name      string            `json:"name"`
	Url       string            `json:"url"`
	ImageName string            `json:"imageName"`
	EnvParam  map[string]string `json:"parameters"`
	JobRuns   []*JobRun         `json:"jobRuns"`
	Config    JobConfiguration
}

// todo test me please
func (j *Job) AppendJobRun(jobRun *JobRun) {
	if j.JobRuns == nil {
		// todo set this default value in another place
		jobRunsSize := j.Config.NbRunBackup
		if jobRunsSize < 1 {
			jobRunsSize = 5
		}
		j.JobRuns = make([]*JobRun, jobRunsSize)
	}
	// todo implements the rolling mecanism
	j.JobRuns[0] = jobRun
}

func (j *Job) GenerateId() error {
	id, err := generateId(j.Id)
	if err == nil {
		j.Id = id
	}
	return err
}

func (j *Job) IsValid() bool {
	if j.Name == "" || j.Url == "" || j.ImageName == "" {
		return false
	}
	return true
}

func (j *Job) String() string {
	return fmt.Sprintf("{Id:%s, Name:%s, Url:%s}", j.Id, j.Name, j.Url)
}

func generateId(id []byte) ([]byte, error) {
	if id != nil && string(id) != "" {
		return nil, errors.New(fmt.Sprintf("the id %+v already defined", string(id)))
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return []byte(strconv.Itoa(r.Int())), nil
}
