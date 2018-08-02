package model

import (
	"fmt"
)

type Job struct {
	Id      []byte    `json:"id"`
	Name    string    `json:"name"`
	GitConf GitConfig `json:"gitConfig"`
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
