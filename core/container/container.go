package container

import (
	"context"

	client "github.com/docker/docker/client"
	"github.com/jeromedoucet/dahu/configuration"
)

// with clientOpts, it will be now simple to override
// option to docker client
var DockerClient ContainerClient = &dockerClient{
	dockerApiVersion: configuration.DockerApiVersion,
	clientOpts:       func(cli *client.Client) error { return nil },
}

type ContainerClient interface {
	CheckRegistryConnection(ctx context.Context, conf RegistryBasicConf) ContainerError
}

type RegistryBasicConf struct {
	Url      string
	User     string
	Password string
}
