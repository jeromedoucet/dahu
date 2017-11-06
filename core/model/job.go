package model

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Configuration detail of
// a Job
type JobConfiguration struct {
	NbRunBackup int `json:"nbRunBackup"` // the number of Run result that are kept
}

type Job struct {
	Id        string            `json:"id"`
	Name      string            `json:"name"`
	Url       string            `json:"url"`
	ImageName string            `json:"imageName"`
	EnvParam  map[string]string `json:"parameters"`
	Config    JobConfiguration
}

func (j *Job) GenerateId() error {
	if j.Id != "" {
		return errors.New(fmt.Sprintf("the id of %+v already defined", j))
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	j.Id = strconv.Itoa(r.Int())

	return nil
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
