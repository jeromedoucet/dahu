package tests

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
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

	waitForGogs()

	return cont.ID
}

func waitForGogs() {
	try := 0
	for {
		if try > 3 {
			panic(errors.New("gogs http port unreachable"))
		}
		conn, err := net.Dial("tcp", "127.0.0.1:10080")
		try++
		if err != nil {
			<-time.After(1 * time.Second)
		} else {
			conn.Close()
			break
		}
	}
}

func StopGogs(id string) {
	client, err := docker.NewClient(endpoint)

	failFast(err)

	rmContainerOpt := docker.RemoveContainerOptions{ID: id, RemoveVolumes: true, Force: true}
	err = client.RemoveContainer(rmContainerOpt)

	failFast(err)
}

// remove the container with the given name
func RemoveContainer(name string) {
	c := exec.Command("docker", []string{"rm", "-f", name}...)
	err := c.Run()
	if err != nil {
		fmt.Println(fmt.Sprintf("got an error %+v when running this command: %s in order to remove a container", err, c.Args))
	}
}

func failFast(err error) {
	if err != nil {
		panic(err)
	}
}
