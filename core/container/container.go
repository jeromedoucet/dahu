package container

import (
	"context"
	"fmt"
	"io"

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

type ContainerEnvs map[string]string

func (e ContainerEnvs) ToArray() []string {
	res := make([]string, len(e), len(e))
	i := 0
	for key, val := range e {
		res[i] = fmt.Sprintf("%s=%s", key, val)
		i++
	}
	return res
}

type ContainerStartConf struct {
	ImageName     string
	RegistryToken string
	Command       []string
	Envs          ContainerEnvs
	ExposedPorts  []Port
	Mounts        []Mount
	WorkingDir    string
	WaitFn        func(ip string) error
}

type ContainerStatus string

const (
	Success  ContainerStatus = "success"
	Error    ContainerStatus = "error"
	Canceled ContainerStatus = "canceled"
)

type ContainerResult struct {
	Status ContainerStatus
	ErrMsg string
}

type ContainerInstance struct {
	Id          string
	Ip          string
	WaitForStop func(chan interface{}) ContainerResult
}

type ContainerRemoveOptions struct {
	RemoveVolumes bool
	Force         bool
}

// encapsulate all operations in containers
type ContainerClient interface {
	CheckRegistryConnection(ctx context.Context, conf RegistryBasicConf) ContainerError
	CreateVolume(ctx context.Context, volumeName string) ContainerError
	RemoveVolume(ctx context.Context, volumeName string) ContainerError
	StartContainer(ctx context.Context, conf ContainerStartConf) (ContainerInstance, ContainerError)
	RemoveContainer(ctx context.Context, id string, options ContainerRemoveOptions) ContainerError
	FollowLogs(ctx context.Context, containerId string, logWriter io.Writer) (ContainerError, chan interface{})
}

type RegistryBasicConf struct {
	Url      string
	User     string
	Password string
}
