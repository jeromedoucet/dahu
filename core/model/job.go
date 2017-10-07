package model

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Job struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (j *Job) GenerateId() error {
	// todo test when already has an id
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	j.Id = strconv.Itoa(r.Int())

	return nil
}

func (j *Job) IsValid() bool {
	if j.Name == "" || j.Url == "" {
		return false
	}
	return true
}

func (j *Job) String() string {
	return fmt.Sprintf("{Id:%s, Name:%s, Url:%s}", j.Id, j.Name, j.Url)
}
