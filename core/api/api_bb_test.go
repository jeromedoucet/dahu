package api_test

import (
	"os"
	"testing"

	"github.com/jeromedoucet/dahu-tests/container"
	"github.com/jeromedoucet/dahu/configuration"
)

var gitRepoIp string

func TestMain(m *testing.M) {
	dockerApiVersion := configuration.DockerApiVersion
	registryId := container.StartDockerRegistry(dockerApiVersion)
	gogsId := container.StartGogs(dockerApiVersion)
	gitRepoDetails := container.FindContainerDetails(gogsId, dockerApiVersion)
	gitRepoIp = gitRepoDetails.Ip
	retCode := m.Run()
	container.StopContainer(gogsId, dockerApiVersion)
	container.StopContainer(registryId, dockerApiVersion)
	os.Exit(retCode)
}
