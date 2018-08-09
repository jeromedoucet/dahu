package model

import (
	"context"

	"github.com/jeromedoucet/dahu/core/container"
)

type DockerRegistryUpdate struct {
	DockerRegistry
	ChangedFields []string
}

type DockerRegistry struct {
	Id                   []byte `json:"id"`
	Name                 string `json:"name"`
	Url                  string `json:"url"`
	User                 string `json:"user"`
	Password             string `json:"password"`
	LastModificationTime int64  `json:"lastModificationTime"`
}

func (r *DockerRegistry) ToPublicModel() {
	r.User = ""
	r.Password = ""
}

func (r DockerRegistry) CheckCredentials(ctx context.Context) container.ContainerError {
	registryConf := container.RegistryBasicConf{User: r.User, Password: r.Password, Url: r.Url}
	return container.DockerClient.CheckRegistryConnection(ctx, registryConf)
}

func (r *DockerRegistry) GenerateId() error {
	id, err := generateId(r.Id)
	if err == nil {
		r.Id = id
	}
	return err
}

func (r *DockerRegistryUpdate) MergeForUpdate(currentRegistry *DockerRegistry) *DockerRegistry {
	res := *currentRegistry
	for _, fieldName := range r.ChangedFields {
		switch fieldName {
		case "Name":
			res.Name = r.Name
		case "Url":
			res.Url = r.Url
		case "User":
			res.User = r.User
		case "Password":
			res.Password = r.Password
		default:
		}
	}
	res.LastModificationTime = r.LastModificationTime
	return &res
}
