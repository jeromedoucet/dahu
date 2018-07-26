package model

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
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

func generateId(id []byte) ([]byte, error) {
	if id != nil && string(id) != "" {
		return nil, errors.New(fmt.Sprintf("the id %+v already defined", string(id)))
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return []byte(strconv.Itoa(r.Int())), nil
}
