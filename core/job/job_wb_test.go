package job

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jeromedoucet/dahu-tests/container"
	"github.com/jeromedoucet/dahu-tests/ssh"
	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
)

// TODO test #FetchSources without auth + with user/pwd + private key
// TODO test #FetchSources with failure on sources volume creation
// TODO test #FetchSources with failure on git clone
// TODO test #FetchSources with failure on log volume creation

// TODO test the websocket in the api_test package

var gitRepoIp string

func TestMain(m *testing.M) {
	dockerApiVersion := configuration.DockerApiVersion
	gogsId := container.StartGogs(dockerApiVersion)
	gitRepoDetails := container.FindContainerDetails(gogsId, dockerApiVersion)
	gitRepoIp = gitRepoDetails.Ip
	retCode := m.Run()
	container.StopContainer(gogsId, dockerApiVersion)
	os.Exit(retCode)
}

// test of FetchSource with ssh private key
// protected by password test.
func TestFetchSourcesWithKeyAuth(t *testing.T) {
	// given
	dockerApiVersion := configuration.DockerApiVersion
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s:10022/tester/test-repo.git", gitRepoIp), Key: ssh.PrivateProtected, KeyPassword: "tester"}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	executionContext := ExecutionContext{BranchName: "master", Context: context.Background(), JobName: "test", ExecutionId: "1"}
	gitVolumeName := fmt.Sprintf("%s-%s-sources", executionContext.JobName, executionContext.ExecutionId)

	// when
	stepExecution := fetchSources(gitConfig, gitVolumeName, executionContext)

	// then
	if stepExecution.Name != "Code fetching" {
		t.Fatalf("expect the step execution name to be 'Code fetching' but is %s", stepExecution.Name)
	}
	if !stepExecution.IsSuccess() {
		t.Fatal("expect the git clone to have been successful, but appear to be failed")
	}
	if !container.VolumeExist(gitVolumeName, dockerApiVersion) {
		t.Fatalf("expect the volume %s to exist, but it doesn't", gitVolumeName)
	}
	// todo check .git
	container.CleanVolume(gitVolumeName, dockerApiVersion)
	// todo check logs
}
