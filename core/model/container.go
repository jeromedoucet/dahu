package model

import (
	"context"

	"github.com/jeromedoucet/dahu/core/container"
)

type DockerRegistry struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (r *DockerRegistry) ToPublicModel() {
	r.User = ""
	r.Password = ""
}

func (r DockerRegistry) CheckCredentials(ctx context.Context) container.ContainerError {
	registryConf := container.RegistryBasicConf{User: r.User, Password: r.Password, Url: r.Url}
	return container.DockerClient.CheckRegistryConnection(ctx, registryConf)
}
