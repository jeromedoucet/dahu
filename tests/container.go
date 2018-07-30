package tests

import (
	"errors"
	"fmt"
	"net"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

const endpoint = "unix:///var/run/docker.sock"

func StartGogs() string {
	client, err := docker.NewClient(endpoint)

	failFast(err)

	exposedPorts := map[docker.Port]struct{}{
		"22/tcp":   {},
		"3000/tcp": {},
	}

	createContConf := docker.Config{
		ExposedPorts: exposedPorts,
		Image:        "jerdct/dahu-gogs",
	}

	portBindings := map[docker.Port][]docker.PortBinding{
		"22/tcp":   {{HostIP: "0.0.0.0", HostPort: "10022"}},
		"3000/tcp": {{HostIP: "0.0.0.0", HostPort: "10080"}},
	}

	createContHostConfig := docker.HostConfig{
		PortBindings:    portBindings,
		PublishAllPorts: true,
		Privileged:      false,
	}

	containerCreationOption := docker.CreateContainerOptions{
		Name:       "gogs_for_test",
		Config:     &createContConf,
		HostConfig: &createContHostConfig,
	}

	var cont *docker.Container
	cont, err = client.CreateContainer(containerCreationOption)

	failFast(err)

	err = client.StartContainer(cont.ID, nil)

	failFast(err)

	waitForService("10080")

	return cont.ID
}

func StartDockerRegistry() string {
	client, err := docker.NewClient(endpoint)

	failFast(err)

	exposedPorts := map[docker.Port]struct{}{
		"5000/tcp": {},
	}

	createContConf := docker.Config{
		ExposedPorts: exposedPorts,
		Image:        "jerdct/dahu-docker-registry",
	}

	portBindings := map[docker.Port][]docker.PortBinding{
		"5000/tcp": {{HostIP: "0.0.0.0", HostPort: "5000"}},
	}

	createContHostConfig := docker.HostConfig{
		PortBindings:    portBindings,
		PublishAllPorts: true,
		Privileged:      false,
	}

	containerCreationOption := docker.CreateContainerOptions{
		Name:       "docker_registry_for_test",
		Config:     &createContConf,
		HostConfig: &createContHostConfig,
	}

	var cont *docker.Container
	cont, err = client.CreateContainer(containerCreationOption)

	failFast(err)

	err = client.StartContainer(cont.ID, nil)

	failFast(err)

	waitForService("5000")

	return cont.ID
}

func StopContainer(id string) {
	client, err := docker.NewClient(endpoint)

	failFast(err)

	rmContainerOpt := docker.RemoveContainerOptions{ID: id, RemoveVolumes: true, Force: true}
	err = client.RemoveContainer(rmContainerOpt)

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
		panic(err)
	}
}
