package container

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	client "github.com/docker/docker/client"
)

type dockerClient struct {
	dockerApiVersion string
	clientOpts       func(*client.Client) error
}

func (d dockerClient) CheckRegistryConnection(ctx context.Context, conf RegistryBasicConf) ContainerError {

	cli, err := client.NewClientWithOpts(client.WithVersion(d.dockerApiVersion), d.clientOpts)
	fmt.Println(err)
	if err != nil {
		return fromDockerToContainerError(err)
	}

	authConfig := types.AuthConfig{
		Username:      conf.User,
		Password:      conf.Password,
		ServerAddress: conf.Url,
	}

	_, err = cli.RegistryLogin(ctx, authConfig)

	return fromDockerToContainerError(err)
}

func fromDockerToContainerError(err error) ContainerError {
	if err == nil {
		return nil
	}
	errStr := err.Error()
	if strings.Contains(errStr, "no such host") {
		return newContainerError(errStr, RegistryNotFound)
	} else if strings.Contains(errStr, "401") {
		return newContainerError(errStr, BadCredentials)
	} else if strings.Contains(errStr, "no basic auth credentials") {
		return newContainerError(errStr, BadCredentials)
	} else {
		return newContainerError(errStr, OtherError)
	}
}
