package scm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/jeromedoucet/dahu-git/types"
	"github.com/jeromedoucet/dahu/core/container"
	"github.com/jeromedoucet/dahu/core/model"
)

// CheckClone is used to check a git repository conf.
// Because this method is made to be used through
// the http api only, it returns a http status directly.
//
// Internally, it Start a dahu-git container and call it
// with given GitConfig. If an error append during
// container handling, 500 http code is returned.
// Else the return code of the call to dahu-git
// is forwarded.
func CheckClone(ctx context.Context, gitConfig model.GitConfig) int {
	dockerCli := container.DockerClient
	dahuGitConf := container.ContainerStartConf{
		ImageName:    "jerdct/dahu-git",
		ExposedPorts: []container.Port{container.Port{Number: "80", Protocol: "tcp"}},
		WaitFn:       waitForDahuGit,
	}
	var dahuGit container.ContainerInstance
	var err error
	dahuGit, err = dockerCli.StartContainer(ctx, dahuGitConf)
	if err != nil {
		return http.StatusInternalServerError
	}

	stopOptions := container.ContainerRemoveOptions{Force: true, RemoveVolumes: true}
	req := buildCheckCloneRequest(gitConfig)
	var reqStatus int
	reqStatus = doCloneForHttp(dahuGit.Ip, req)

	// for the moment, the choice is make
	// to consider that the result of the container
	// stop must be show to the user as a system error.
	// this is because if we have trouble to stop a simple
	// container, then the whole system is in deep trouble.
	// Meaning there is a need for deeper investigations.
	err = dockerCli.RemoveContainer(ctx, dahuGit.Id, stopOptions)
	if err != nil {
		return http.StatusInternalServerError
	}

	return reqStatus
}

func buildCheckCloneRequest(g model.GitConfig) types.CloneRequest {
	req := types.CloneRequest{Branch: "master", NoCheckout: true}
	if g.HttpAuth != nil {
		req.UseHttp = true
		req.HttpAuth = types.HttpAuth{
			Url:      g.HttpAuth.Url,
			User:     g.HttpAuth.User,
			Password: g.HttpAuth.Password,
		}
	} else if g.SshAuth != nil {
		req.UseSsh = true
		req.SshAuth = types.SshAuth{
			Url:         g.SshAuth.Url,
			Key:         g.SshAuth.Key,
			KeyPassword: g.SshAuth.KeyPassword,
		}
	}
	return req
}

type CloneConfiguration struct {
	GitConfig  model.GitConfig
	BranchName string
	VolumeName string
	LogWriter  io.Writer
	NetworkId  string
}

// todo factorisation with CheckClone()
func Clone(ctx context.Context, conf CloneConfiguration) error {
	dockerCli := container.DockerClient
	var err error
	destinationFolder := "/data"

	err = dockerCli.CreateVolume(ctx, conf.VolumeName)
	if err != nil {
		return err
	}

	mounts := []container.Mount{container.Mount{Source: conf.VolumeName, Destination: destinationFolder}}
	dahuGitConf := container.ContainerStartConf{
		ImageName:    "jerdct/dahu-git",
		ExposedPorts: []container.Port{container.Port{Number: "80", Protocol: "tcp"}},
		WaitFn:       waitForDahuGit,
		Mounts:       mounts,
		NetworkId:    conf.NetworkId,
	}

	var dahuGit container.ContainerInstance
	dahuGit, err = dockerCli.StartContainer(ctx, dahuGitConf)

	if err != nil {
		return err
	}

	var waitLog chan interface{}
	err, waitLog = dockerCli.FollowLogs(ctx, dahuGit.Id, conf.LogWriter)

	if err != nil {
		return err
	}

	removeOptions := container.ContainerRemoveOptions{Force: true, RemoveVolumes: true}

	req := buildCloneRequest(conf)
	err = doCloneForJob(dahuGit.Ip, req)
	if err != nil {
		// Don't forget to stop the container anyway !
		dockerCli.RemoveContainer(ctx, dahuGit.Id, removeOptions)
		<-waitLog
		return err
	}

	// for the moment, the choice is make
	// to consider that the result of the container
	// stop must be show to the user as a system error.
	// this is because if we have trouble to stop a simple
	// container, then the whole system is in deep trouble.
	// Meaning there is a need for deeper investigations.
	err = dockerCli.RemoveContainer(ctx, dahuGit.Id, removeOptions)
	<-waitLog

	return err
}

func buildCloneRequest(conf CloneConfiguration) types.CloneRequest {
	req := types.CloneRequest{Branch: conf.BranchName}
	if conf.GitConfig.HttpAuth != nil {
		req.UseHttp = true
		req.HttpAuth = types.HttpAuth{
			Url:      conf.GitConfig.HttpAuth.Url,
			User:     conf.GitConfig.HttpAuth.User,
			Password: conf.GitConfig.HttpAuth.Password,
		}
	} else if conf.GitConfig.SshAuth != nil {
		req.UseSsh = true
		req.SshAuth = types.SshAuth{
			Url:         conf.GitConfig.SshAuth.Url,
			Key:         conf.GitConfig.SshAuth.Key,
			KeyPassword: conf.GitConfig.SshAuth.KeyPassword,
		}
	}
	return req
}

func waitForDahuGit(ip string) error {
	try := 0
	for {
		if try > 20 {
			return errors.New("ERROR => dahu-git unreachable")
		}
		resp, err := http.Get(fmt.Sprintf("http://%s/status", ip))
		try++
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("Failed to reach dah-git at %d attempt", try)
			<-time.After(1 * time.Second)
		} else {
			return nil
		}
	}
}

func doCloneForJob(ip string, req types.CloneRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/clone", ip), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return errors.New("ERROR >> Error cloning the repo")
	}
	return nil
}

func doCloneForHttp(ip string, req types.CloneRequest) int {
	defaultStatus := http.StatusInternalServerError
	body, err := json.Marshal(req)
	if err != nil {
		return defaultStatus
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/clone", ip), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return defaultStatus
	}
	return resp.StatusCode
}
