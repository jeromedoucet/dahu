package tests

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jeromedoucet/dahu/configuration"
)

func StartGogs() string {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.WithVersion(configuration.DockerApiVersion))

	failFast(err)

	internalSshPort, _ := nat.NewPort("tcp", "22")
	internalHttpPort, _ := nat.NewPort("tcp", "3000")

	exposedPorts := nat.PortSet{
		internalSshPort:  {},
		internalHttpPort: {},
	}

	containerConf := &container.Config{Image: "jerdct/dahu-gogs", ExposedPorts: exposedPorts}

	externalSshPort := nat.PortBinding{HostIP: "0.0.0.0", HostPort: "10022"}
	externalHttpPort := nat.PortBinding{HostIP: "0.0.0.0", HostPort: "10080"}

	portBindings := nat.PortMap{
		internalSshPort:  []nat.PortBinding{externalSshPort},
		internalHttpPort: []nat.PortBinding{externalHttpPort},
	}

	hostConfig := &container.HostConfig{PortBindings: portBindings}

	networkConfig := &network.NetworkingConfig{}

	var createdContainer container.ContainerCreateCreatedBody
	createdContainer, err = cli.ContainerCreate(ctx, containerConf, hostConfig, networkConfig, "gogs_for_test")

	failFast(err)

	err = cli.ContainerStart(ctx, createdContainer.ID, types.ContainerStartOptions{})

	failFast(err)

	waitForService("10080")

	return createdContainer.ID
}

func StartDockerRegistry() string {

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.WithVersion(configuration.DockerApiVersion))

	failFast(err)

	internalPort, _ := nat.NewPort("tcp", "5000")

	exposedPorts := nat.PortSet{
		internalPort: {},
	}

	containerConf := &container.Config{Image: "jerdct/dahu-docker-registry", ExposedPorts: exposedPorts}

	externalPort := nat.PortBinding{HostIP: "0.0.0.0", HostPort: "5000"}

	portBindings := nat.PortMap{
		internalPort: []nat.PortBinding{externalPort},
	}

	hostConfig := &container.HostConfig{PortBindings: portBindings}

	networkConfig := &network.NetworkingConfig{}

	var createdContainer container.ContainerCreateCreatedBody
	createdContainer, err = cli.ContainerCreate(ctx, containerConf, hostConfig, networkConfig, "docker_registry_for_test")

	failFast(err)

	err = cli.ContainerStart(ctx, createdContainer.ID, types.ContainerStartOptions{})

	failFast(err)

	waitForService("5000")

	return createdContainer.ID
}

func StopContainer(id string) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.WithVersion(configuration.DockerApiVersion))

	failFast(err)

	removeOpt := types.ContainerRemoveOptions{Force: true}

	err = cli.ContainerRemove(ctx, id, removeOpt)

	failFast(err)
}

func waitForService(tcpPort string) {
	try := 0
	for {
		if try > 3 {
			panic(errors.New("gogs http port unreachable"))
		}
		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", tcpPort))
		try++
		if err != nil {
			<-time.After(1 * time.Second)
		} else {
			conn.Close()
			break
		}
	}
}

func failFast(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
