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

type Mount struct {
	Source      string
	Destination string
}

type Port struct {
	Number   string
	Protocol string
}

type ContainerStartConf struct {
	ImageName    string
	ExposedPorts []Port
	Mounts       []Mount
	WaitFn       func() error
}

type ContainerInstance struct {
	Id string
	Ip string
}

type ContainerStopOptions struct {
	RemoveVolumes bool
	Force         bool
}

// encapsulate all operations in containers
type ContainerClient interface {
	CheckRegistryConnection(ctx context.Context, conf RegistryBasicConf) ContainerError
	CreateVolume(ctx context.Context, volumeName string) ContainerError
	StartContainer(ctx context.Context, conf ContainerStartConf) (ContainerInstance, ContainerError)
	StopContainer(ctx context.Context, id string, options ContainerStopOptions) ContainerError
}

type RegistryBasicConf struct {
	Url      string
	User     string
	Password string
}
