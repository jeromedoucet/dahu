package model

import (
	"context"
	"strconv"
	"time"

	"github.com/jeromedoucet/dahu/core/container"
)

type DockerRegistryUpdate struct {
	DockerRegistry
	ChangedFields []string `json:"changedFields"`
}

type DockerRegistry struct {
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	Url                  string `json:"url"`
	User                 string `json:"user"`
	Password             string `json:"password"`
	LastModificationTime string `json:"lastModificationTime"`
}

func (r *DockerRegistry) ToPublicModel() {
	r.Password = ""
}

func (r DockerRegistry) CheckCredentials(ctx context.Context) container.ContainerError {
	registryConf := container.RegistryBasicConf{User: r.User, Password: r.Password, Url: r.Url}
	return container.DockerClient.CheckRegistryConnection(ctx, registryConf)
}

func (r *DockerRegistry) GenerateId() error {
	id, err := generateId([]byte(r.Id))
	if err == nil {
		r.Id = string(id)
	}
	return err
}

// update the LastModificationTimeField
func (r *DockerRegistry) NewLastModificationTime() {
	timeStamp := time.Now().UnixNano()
	r.LastModificationTime = strconv.Itoa(int(timeStamp))
}

func (r *DockerRegistryUpdate) MergeForUpdate(currentRegistry *DockerRegistry) *DockerRegistry {
	res := *currentRegistry
	for _, fieldName := range r.ChangedFields {
		switch fieldName {
		case "name":
			res.Name = r.Name
		case "url":
			res.Url = r.Url
		case "user":
			res.User = r.User
		case "password":
			res.Password = r.Password
		default:
		}
	}
	res.LastModificationTime = r.LastModificationTime
	return &res
}
