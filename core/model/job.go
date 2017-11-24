package model

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strconv"
	"time"
)

type RunStatus int

// available run status
const (
	CREATED RunStatus = 1 + iota
	RUNNING
	CANCELED
	SUCCESS
	FAILURE
)

func isAvailableRunStatus(s RunStatus) bool {
	return !(s < CREATED || s > FAILURE)
}

// use by request that
// will start a job run.
type RunRequest struct {
	OpenWs bool `json:"openWs"`
}

// A job run store informations on a run
// for a given Job
type JobRun struct {
	Id            []byte     `json:"id"`
	ContainerName string     `json:"containerName"`
	Status        RunStatus  `json:"runStatus"`
	StartTime     *time.Time `json:"startTime"`
	EndTime       *time.Time `json:"endTime"`
	Version       int64      `json: "version"`
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

// return true if the JobRun instance
// has enought information to be registered
// it should have a non nil ContainerName and
// a RunStatus
func (j *JobRun) IsValid() bool {
	return j.ContainerName != "" && isAvailableRunStatus(j.Status)
}

// Configuration detail of
// a Job
type JobConfiguration struct {
	NbRunBackup int `json:"nbRunBackup"` // the number of Run result that are kept
}

type Job struct {
	Id        []byte            `json:"id"`
	Name      string            `json:"name"`
	Url       string            `json:"url"`
	ImageName string            `json:"imageName"`
	EnvParam  map[string]string `json:"parameters"`
	JobRuns   []*JobRun         `json:"jobRuns"`
	Config    JobConfiguration
}

// Append with rolling policy a
// JobRun on this Job
// If nil slice, create a slice
// on 5 size
func (j *Job) AppendJobRun(jobRun *JobRun) {
	if j.JobRuns == nil {
		jobRunsSize := j.Config.NbRunBackup
		if jobRunsSize < 1 {
			jobRunsSize = 5
		}
		j.JobRuns = make([]*JobRun, jobRunsSize)
	}
	rollingJr := jobRun
	for i, v := range j.JobRuns {
		j.JobRuns[i] = rollingJr
		rollingJr = v
	}
	// when a joRun is removed from a job
	// the underlying container which does still
	// exist to keep logs should now
	// be remove.
	if rollingJr != nil {
		cmd := exec.Command("docker", "rm", "-f", rollingJr.ContainerName)
		err := cmd.Run()
		if err != nil {
			log.Printf("WARN >> Encounter error when trying to remove container %s : %+v", rollingJr.ContainerName, err)
		}
	}
}

func (j *Job) FindJobRun(id []byte) (*JobRun, error) {
	for _, v := range j.JobRuns {
		if v != nil && string(v.Id) == string(id) {
			return v, nil
		}
	}
	return nil, NewNoMorePersisted(fmt.Sprintf("WARN >> Job %s, the JobRun %s is not persisted anymore.", string(j.Id), string(id)))
}

func (j *Job) UpdateJobRun(jobRun *JobRun) *JobRun {
	existingJobRun, err := j.FindJobRun(jobRun.Id)
	if err != nil {
		return jobRun
	}
	if jobRun.Version != existingJobRun.Version {
		return existingJobRun
	}
	jobRun.Version = time.Now().UnixNano()
	for i, v := range j.JobRuns {
		if v != nil && string(v.Id) == string(jobRun.Id) {
			j.JobRuns[i] = jobRun
			break
		}
	}
	return jobRun
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
