package api_test

import (
	"os"
	"testing"

	"github.com/jeromedoucet/dahu-tests/container"
	"github.com/jeromedoucet/dahu/configuration"
)

func TestMain(m *testing.M) {
	dockerApiVersion := configuration.DockerApiVersion
	gogsId := container.StartGogs(dockerApiVersion)
	registryId := container.StartDockerRegistry(dockerApiVersion)
	gitRepoDetails := container.FindContainerDetails(gogsId, dockerApiVersion)
	gitRepoIp = gitRepoDetails.Ip
	retCode := m.Run()
	container.StopContainer(gogsId, dockerApiVersion)
	container.StopContainer(registryId, dockerApiVersion)
	os.Exit(retCode)
}
