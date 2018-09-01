package job

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jeromedoucet/dahu-tests/container"
	"github.com/jeromedoucet/dahu-tests/ssh"
	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
)

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
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/test-repo.git", gitRepoIp), Key: ssh.PrivateProtected, KeyPassword: "tester"}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	executionContext := ExecutionContext{BranchName: "master", Context: context.Background(), JobName: "test", ExecutionId: "1"}
	gitVolumeName := fmt.Sprintf("%s-%s-sources", executionContext.JobName, executionContext.ExecutionId)
	job := model.Job{GitConf: gitConfig}
	exec := execution{executionContext: executionContext, job: job, sourcesVolume: gitVolumeName}

	// when
	stepExecution := exec.fetchSources()

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
	container.CleanVolume(gitVolumeName, dockerApiVersion)
	if !strings.Contains(stepExecution.Logs, "Clone finished without error") {
		t.Fatalf("expect the logs of the clone to contains 'Clone finished without error', but got %s", stepExecution.Logs)
	}
}

// test of FetchSource with one error during clone (authentication)
func TestFetchSourcesWithBadAuth(t *testing.T) {
	// given
	dockerApiVersion := configuration.DockerApiVersion
	authConfig := model.SshAuthConfig{Url: fmt.Sprintf("ssh://git@%s/tester/test-repo.git", gitRepoIp)}
	gitConfig := model.GitConfig{SshAuth: &authConfig}
	executionContext := ExecutionContext{BranchName: "master", Context: context.Background(), JobName: "test", ExecutionId: "1"}
	gitVolumeName := fmt.Sprintf("%s-%s-sources", executionContext.JobName, executionContext.ExecutionId)
	job := model.Job{GitConf: gitConfig}
	exec := execution{executionContext: executionContext, job: job, sourcesVolume: gitVolumeName}

	// when
	stepExecution := exec.fetchSources()

	// then
	if stepExecution.Name != "Code fetching" {
		t.Fatalf("expect the step execution name to be 'Code fetching' but is %s", stepExecution.Name)
	}
	if stepExecution.IsSuccess() {
		t.Fatal("expect the git clone to have failed, but appear to have succed")
	}
	if !container.VolumeExist(gitVolumeName, dockerApiVersion) {
		t.Fatalf("expect the volume %s to exist, but it doesn't", gitVolumeName)
	}
	container.CleanVolume(gitVolumeName, dockerApiVersion)
}

// todo test container creation failure
