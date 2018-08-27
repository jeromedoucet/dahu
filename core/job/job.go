package job

import (
	"context"

	"github.com/jeromedoucet/dahu/core/container"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/scm"
)

// how it should work
// 1. Job submit :
//     - one job structure is passed to 'Start' function
//     - if dahu is closing, an error is returned (through recovering a panic cause by writting in a closed channel)
//     - if not, an execution start (a goroutine - with context ? - is started), and nil is returned

// 2. Job execution :
//     - initialization of jobExecution instance with default instance and save. Must add an update with jobId and JobExecution (with an id) and don't erase running one
//     A : git clone
//       - update the jobExecution for the first step status => running
//       - git clone the correct branch and target a docker volume. Stdout of clone should be saved and update the JobExecution with right status (success or failed) and the
//       volume name
//       - if failed, stop.
//       - else go to next

// note on job execution persistence. Must keep the volume + the logs + the status of each step

// note on gracefull shutdown. TODO

// note on configuration : websocket should be passed as an optional entry, branch as a mandatory one

// note on scm pull => include it as the first step

// note => handle the context.Context corretly

// this contains all contextual
// resources of the build. Branch Name,
// webSocket connection for 'live' notifications
// or the event that trigger the build
//
// some of this information are lately
// saved into the model (for exemple the branch bane)
type ExecutionContext struct {
	BranchName  string
	Context     context.Context // todo test non-nil ?
	JobName     string          // TODO clean the string somewhere to match docker volume name requirement
	ExecutionId string
}

func Start(job model.Job) error {
	return nil
}

// todo, pass the gitVolumeName
func fetchSources(repoConf model.GitConfig, sourcesVolume string, executionContext ExecutionContext) (stepExecution model.StepExecution) {
	//var err error
	containerCli := container.DockerClient

	containerCli.CreateVolume(executionContext.Context, sourcesVolume) // todo handle error

	cloneConf := scm.CloneConfiguration{
		GitConfig:  repoConf,
		BranchName: executionContext.BranchName,
		VolumeName: sourcesVolume,
	}

	// todo pass something to pipe the logs
	// todo think of persistence, and log persistence.
	scm.Clone(executionContext.Context, cloneConf) // todo handle error
	return model.StepExecution{Name: "Code fetching", Status: model.Success}
}
