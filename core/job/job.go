// job package is where job execution logic is. Such execution has step. Basically,
// step are a logical split of the process, for instance "run tests" or "deploy build".
// The first step is implicit and alway exist: fetching the repository from svc server.
//
// - SVC step
// The code handling this process is not included in that package. A dedicated container is used
// to perform the job. This allow to have less code in the dahu core and to extend capabilities easily.
// Such pattern is heavily used in Dahu. The local copy of sources is stored in a docker volume that will
// be available for the next steps.
//
// - User-defined steps
// This is the business logic of the job. typically, you may have one step to fetch dependencies,
// then one or two steps for the tests and finally build deployment and notifications.
// To run a step, Dahu will pull an image (from a public or a private repository), start if needed some
// required services (see bellow), and start a container from the pulled images with a given command
// and some properties and configuration options. This container is executed in a dedicated network.
// This is required for services (see bellow).
// If the step stop successfully, the next step will be started. if not, the job stop.
//
// - Services
// For some kind of steps (integration test for example), some running process are needed. Dahu has the concept of
// service to fill that need. A service is a dependency of a step and consist of a docker image with some configuration options.
// It is run inside the same network than the related step and is not reachable from outside. Services are accessible through there
// names from the step container.
//
// - Cancelation
// Steps can be canceled anytime. To achieve that, there is an internal scheduler keeping a reference to a channel for all job execution process.
// When an execution start, it is registered on that scheduler. The unregistration is done at the end of the execution, regardless of the result.
//
// - Notifications
// A websocket channel is available to listen for job executions's events. An event may be a job start, a job stop, a step start, a step stop or a log event.

package job

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/container"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
	"github.com/jeromedoucet/dahu/core/scm"
)

// Start launch a new job execution. It runs in a dedicated goroutine.
func Start(job model.Job, branchName string, conf *configuration.Conf, ctx context.Context) model.JobExecution {
	jobExecution := model.JobExecution{BranchName: branchName}
	jobExecution.GenerateId()
	e := execution{
		job:           job,
		jobExecution:  jobExecution,
		ctx:           ctx,
		sourcesVolume: fmt.Sprintf("%s-%s-sources", job.Name, jobExecution.Id),
		conf:          conf,
		repository:    persistence.GetRepository(conf),
	}
	go e.run()
	return jobExecution
}

// Inner type that contains all informations
// on a job execution. The job itself, its context,
// and some related information like the name of
// the volume where the sources are.
type execution struct {
	cancelChan    chan interface{}
	job           model.Job
	jobExecution  model.JobExecution
	ctx           context.Context
	sourcesVolume string
	conf          *configuration.Conf
	repository    persistence.Repository
	networkId     string
}

// Run contains the main loop of a job execution process.
// it start each step of an execution and handle their result.
func (e execution) run() {

	e.cancelChan = registerJobExecution(string(e.job.Id), e.jobExecution.Id)

	Broadcast(string(e.job.Id), model.Event{
		Type:        model.JobStart,
		ExecutionId: e.jobExecution.Id,
		Value:       fmt.Sprintf("Start execute job %s on branch %s", e.job.Name, e.jobExecution.BranchName),
	})

	// all container are attached to a custom network
	err, networkId := container.DockerClient.CreateNetwork(e.ctx, fmt.Sprintf("network-%s", e.jobExecution.Id))
	if err != nil {
		// todo update endJob to accept an optional error
		Broadcast(string(e.job.Id), model.Event{
			Type:        model.NewLog,
			ExecutionId: e.jobExecution.Id,
			Value:       fmt.Sprintf("Error when creating a network : %s ", err.Error()),
		})
		e.endJob(model.Failure)
		return
	}

	e.networkId = networkId

	fetchExecution := &model.StepExecution{Name: "Code fetching", Status: model.Running}
	e.jobExecution.Steps = append(e.jobExecution.Steps, fetchExecution)

	e.repository.UpsertJobExecution(e.ctx, string(e.job.Id), &e.jobExecution)
	e.fetchSources(fetchExecution)
	e.repository.UpsertJobExecution(e.ctx, string(e.job.Id), &e.jobExecution)

	if fetchExecution.Status == model.Success {
		for _, step := range e.job.Steps {

			stepExecution := &model.StepExecution{Name: step.Name, Status: model.Running}
			e.jobExecution.Steps = append(e.jobExecution.Steps, stepExecution)

			e.repository.UpsertJobExecution(e.ctx, string(e.job.Id), &e.jobExecution)
			e.executeStep(&step, stepExecution)
			e.repository.UpsertJobExecution(e.ctx, string(e.job.Id), &e.jobExecution)

			if stepExecution.Status == model.Failure || stepExecution.Status == model.Canceled {
				e.endJob(stepExecution.Status)
				return
			}
		}
	} else {
		e.endJob(model.Failure)
		return
	}
	e.endJob(model.Success)
}

// endJob handle terminal operation of a job execution. Workspace
// cleaning, end event broacasting, execution unregistrations...
func (e execution) endJob(terminationStatus model.ExecutionStatus) {
	containerCli := container.DockerClient
	if e.job.RemoveWorkspace {
		err := containerCli.RemoveVolume(e.ctx, e.sourcesVolume) // TODO handle that error a clean way
		if err != nil {
			log.Printf("ERROR >> run encounter error : %s", err.Error())
		}
	}
	if terminationStatus == model.Success {
		Broadcast(string(e.job.Id), model.Event{
			Type:        model.JobSucceed,
			ExecutionId: e.jobExecution.Id,
			Value:       fmt.Sprintf("Finished job %s execution on branch %s", e.job.Name, e.jobExecution.BranchName),
		})
	} else if terminationStatus == model.Failure {
		Broadcast(string(e.job.Id), model.Event{
			Type:        model.JobFailed,
			ExecutionId: e.jobExecution.Id,
			Value:       fmt.Sprintf("Finished job %s execution on branch %s with failure", e.job.Name, e.jobExecution.BranchName),
		})
	} else {
		Broadcast(string(e.job.Id), model.Event{
			Type:        model.JobCanceled,
			ExecutionId: e.jobExecution.Id,
			Value:       fmt.Sprintf("Canceled job %s execution on branch %s", e.job.Name, e.jobExecution.BranchName),
		})
	}

	// Don't forget that. This is permit to clean references
	// in the job execution scheduler.
	unRegisterJobExecution(string(e.job.Id), e.jobExecution.Id)
	e.repository.UpsertJobExecution(e.ctx, string(e.job.Id), &e.jobExecution)

	// at the end, the network should be remove
	containerCli.DeleteNetwork(e.ctx, e.networkId)
}

// fetchSources is the first step of a job execution. Like
// its mame suggests, it will get the sources.
func (e execution) fetchSources(stepExecution *model.StepExecution) {

	Broadcast(string(e.job.Id), model.Event{
		Type:        model.StepStart,
		ExecutionId: e.jobExecution.Id,
		Value:       "Start fetching code",
	})

	containerCli := container.DockerClient

	containerCli.CreateVolume(e.ctx, e.sourcesVolume) // TODO handle error

	w := &logWriter{
		jobId:       string(e.job.Id),
		executionId: e.jobExecution.Id,
	}

	cloneConf := scm.CloneConfiguration{
		GitConfig:  e.job.GitConf,
		BranchName: e.jobExecution.BranchName,
		VolumeName: e.sourcesVolume,
		LogWriter:  w,
		NetworkId:  e.networkId,
	}

	err := scm.Clone(e.ctx, cloneConf)
	if err == nil {
		stepExecution.Status = model.Success
		Broadcast(string(e.job.Id), model.Event{
			Type:        model.StepSucceed,
			ExecutionId: e.jobExecution.Id,
			Value:       "Succeed fetching code",
		})
	} else {
		log.Printf("Job >> issue when fetching sources %s", err.Error())
		stepExecution.Status = model.Failure
		Broadcast(string(e.job.Id), model.Event{
			Type:        model.StepFailed,
			ExecutionId: e.jobExecution.Id,
			Value:       "Failed fetching code",
		})
	}
	stepExecution.Logs = string(w.logs)
}

// executeStep is responsible for preparing a step, run it, notifying
// events.
// It heavily rely on container package, that abstract container manipulations.
func (e execution) executeStep(step *model.Step, stepExecution *model.StepExecution) {
	var c container.ContainerInstance
	var err error
	var services []*container.ContainerInstance

	Broadcast(string(e.job.Id), model.Event{
		Type:        model.StepStart,
		ExecutionId: e.jobExecution.Id,
		Value:       fmt.Sprintf("Start %s", step.Name),
	})

	err, services = e.startServices(step)
	defer e.stopServices(services)

	if err != nil {
		Broadcast(string(e.job.Id), model.Event{
			Type:        model.StepFailed,
			ExecutionId: e.jobExecution.Id,
			Value:       fmt.Sprintf("%s has failed : %s", step.Name, err.Error()),
		})
		stepExecution.Status = model.Failure
		stepExecution.Logs = err.Error()
		return
	}

	registryToken := getRegistryAuth(step.Image)

	dockerCli := container.DockerClient
	mounts := []container.Mount{container.Mount{Source: e.sourcesVolume, Destination: step.MountingPoint}}
	stepConf := container.ContainerStartConf{
		ImageName:     step.Image.Name,
		RegistryToken: registryToken,
		Mounts:        mounts,
		Command:       step.Command,
		WorkingDir:    step.MountingPoint,
		Envs:          step.ComputeEnvs(),
		NetworkId:     e.networkId,
	}

	c, err = dockerCli.StartContainer(e.ctx, stepConf)

	if err != nil {
		Broadcast(string(e.job.Id), model.Event{
			Type:        model.StepFailed,
			ExecutionId: e.jobExecution.Id,
			Value:       fmt.Sprintf("%s has failed : %s", step.Name, err.Error()),
		})
		stepExecution.Status = model.Failure
		stepExecution.Logs = err.Error()
		return
	}

	w := &logWriter{
		jobId:       string(e.job.Id),
		executionId: e.jobExecution.Id,
	}

	err, _ = dockerCli.FollowLogs(e.ctx, c.Id, w)

	if err != nil {
		Broadcast(string(e.job.Id), model.Event{
			Type:        model.StepFailed,
			ExecutionId: e.jobExecution.Id,
			Value:       fmt.Sprintf("%s failed : %s", step.Name, err.Error()),
		})
		stepExecution.Status = model.Failure
		stepExecution.Logs = err.Error()
		return
	}

	containerResult := c.WaitForStop(e.cancelChan)

	if containerResult.Status == container.Success {
		stepExecution.Status = model.Success
		stepExecution.Logs = string(w.logs)
		Broadcast(string(e.job.Id), model.Event{
			Type:        model.StepSucceed,
			ExecutionId: e.jobExecution.Id,
			Value:       fmt.Sprintf("Finished %s", step.Name),
		})
	} else if containerResult.Status == container.Error {
		stepExecution.Status = model.Failure
		stepExecution.Logs = string(w.logs)
		Broadcast(string(e.job.Id), model.Event{
			Type:        model.StepFailed,
			ExecutionId: e.jobExecution.Id,
			Value:       containerResult.ErrMsg,
		})
	} else {
		stepExecution.Status = model.Canceled
		stepExecution.Logs = string(w.logs)
		Broadcast(string(e.job.Id), model.Event{
			Type:        model.StepCanceled,
			ExecutionId: e.jobExecution.Id,
			Value:       fmt.Sprintf("Finished %s", step.Name),
		})
	}

	removeOptions := container.ContainerRemoveOptions{Force: true, RemoveVolumes: true}

	dockerCli.RemoveContainer(e.ctx, c.Id, removeOptions)
}

// startServices launch all services registered for a
// step. It return an array of containerInstances and
// an error if something went wrong during services launch.
func (e execution) startServices(step *model.Step) (error, []*container.ContainerInstance) {

	res := make([]*container.ContainerInstance, len(step.Services), len(step.Services))
	dockerCli := container.DockerClient
	for i, service := range step.Services {
		exposedPorts := make([]container.Port, len(service.ExposedPorts), len(service.ExposedPorts))
		for j, port := range service.ExposedPorts {
			exposedPorts[j] = container.Port{Number: strconv.Itoa(port.Num), Protocol: port.Prototype}
		}
		registryToken := getRegistryAuth(service.Image)
		serviceConf := container.ContainerStartConf{
			ContainerName: service.Name,
			ImageName:     service.Image.Name,
			RegistryToken: registryToken,
			ExposedPorts:  exposedPorts,
			NetworkId:     e.networkId,
		}
		c, err := dockerCli.StartContainer(e.ctx, serviceConf)
		if err != nil {
			return err, res
		} else {
			res[i] = &c
		}
	}

	return nil, res
}

func (e execution) stopServices(servicesInstance []*container.ContainerInstance) error {
	dockerCli := container.DockerClient
	for _, instance := range servicesInstance {
		removeConf := container.ContainerRemoveOptions{true, true}
		err := dockerCli.RemoveContainer(e.ctx, instance.Id, removeConf)
		if err != nil {
			return err
		}
	}
	return nil
}

func getRegistryAuth(image model.Image) string {
	if image.Registry != nil {
		registryAuth := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: image.Registry.User,
			Password: image.Registry.Password,
		}
		data, _ := json.Marshal(registryAuth)
		fmt.Println(string(data))
		return base64.StdEncoding.EncodeToString(data)
	} else {
		return ""
	}
}

type logWriter struct {
	jobId       string
	executionId string
	logs        []byte
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	if len(p) > 0 {
		Broadcast(string(l.jobId), model.Event{
			Type:        model.NewLog,
			ExecutionId: l.executionId,
			Value:       strings.TrimSpace(string(p)),
		})

		l.logs = append(l.logs, p...)
	}
	return len(p), nil
}
