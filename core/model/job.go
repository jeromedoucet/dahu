package model

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Job struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Url        string   `json:"url"`
	ImageName  string   `json:"imageName"`
	Parameters []string `json:"parameters"`
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
