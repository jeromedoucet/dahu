package container

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	client "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type dockerClient struct {
	dockerApiVersion string
	clientOpts       func(*client.Client) error
}

func (d dockerClient) CheckRegistryConnection(ctx context.Context, conf RegistryBasicConf) ContainerError {

	cli, err := client.NewClientWithOpts(client.WithVersion(d.dockerApiVersion), d.clientOpts)
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

func (d dockerClient) StartContainer(ctx context.Context, conf ContainerStartConf) (ContainerInstance, ContainerError) {
	instance := ContainerInstance{}
	cli, err := client.NewClientWithOpts(client.WithVersion(d.dockerApiVersion))
	if err != nil {
		return instance, fromDockerToContainerError(err)
	}
	defer cli.Close() // todo consider handling the error

	// pull the image right now
	err = pullImage(ctx, conf.ImageName, cli)
	if err != nil {
		return instance, fromDockerToContainerError(err)
	}

	// container configuration
	exposedPorts, err := createPortsConf(conf.ExposedPorts)
	if err != nil {
		return instance, fromDockerToContainerError(err)
	}
	mounts := createMounts(conf.Mounts)
	containerConf := &container.Config{Image: conf.ImageName, ExposedPorts: exposedPorts}
	hostConfig := &container.HostConfig{Mounts: mounts}
	networkConfig := &network.NetworkingConfig{}

	// the first step is to create the container. It is worth to notice
	// that it is not running yet
	var createdContainer container.ContainerCreateCreatedBody
	createdContainer, err = cli.ContainerCreate(ctx, containerConf, hostConfig, networkConfig, "")
	if err != nil {
		return instance, fromDockerToContainerError(err)
	}

	// Now the container will start
	err = cli.ContainerStart(ctx, createdContainer.ID, types.ContainerStartOptions{})
	if err != nil {
		return instance, fromDockerToContainerError(err)
	}

	// TODO at this step, when an error is raised, don't forget
	// to STOP the container if someting bad append

	// but we still don't have its ip, which
	// is quite important is some case
	var inspectResult types.ContainerJSON
	inspectResult, err = cli.ContainerInspect(ctx, createdContainer.ID)
	if err != nil {
		return instance, fromDockerToContainerError(err)
	}

	// If there is a specific function that must be used
	// to check if the container is ready, it must be executed now.
	if conf.WaitFn != nil {
		// TODO set a time out here !
		err = conf.WaitFn()
		if err != nil {
			return instance, fromDockerToContainerError(err)
		}
	}

	// everything fine. fill the instance structure
	instance.Id = createdContainer.ID
	instance.Ip = inspectResult.NetworkSettings.IPAddress
	return instance, nil
}

func (d dockerClient) StopContainer(ctx context.Context, id string, options ContainerStopOptions) ContainerError {
	cli, err := client.NewClientWithOpts(client.WithVersion(d.dockerApiVersion))
	if err != nil {
		return fromDockerToContainerError(err)
	}
	defer cli.Close()

	removeOpt := types.ContainerRemoveOptions{Force: options.Force, RemoveVolumes: options.RemoveVolumes}

	return fromDockerToContainerError(cli.ContainerRemove(ctx, id, removeOpt))
}

func (d dockerClient) CreateVolume(ctx context.Context, volumeName string) ContainerError {
	cli, err := client.NewClientWithOpts(client.WithVersion(d.dockerApiVersion))
	if err != nil {
		return fromDockerToContainerError(err)
	}
	defer cli.Close()

	_, err = cli.VolumeCreate(ctx, volume.VolumeCreateBody{Name: volumeName})
	return fromDockerToContainerError(err)
}

func pullImage(ctx context.Context, imageName string, cli *client.Client) error {
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return err
	}

	return out.Close()
}

type empty struct{}

func createPortsConf(exposedPorts []Port) (nat.PortSet, ContainerError) {
	res := nat.PortSet{}
	for _, port := range exposedPorts {
		exposedPort, err := nat.NewPort(port.Protocol, port.Number)
		if err != nil {
			return res, fromDockerToContainerError(err)
		}
		res[exposedPort] = empty{}
	}
	return res, nil
}

func createMounts(mounts []Mount) []mount.Mount {
	res := []mount.Mount{}
	for _, m := range mounts {
		res = append(res, mount.Mount{Type: mount.TypeVolume, Source: m.Source, Target: m.Destination})
	}
	return res
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
