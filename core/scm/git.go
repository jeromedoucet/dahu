package scm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
		ImageName:    "jerdct/dahu-git", // TODO handle tag
		ExposedPorts: []container.Port{container.Port{Number: "80", Protocol: "tcp"}},
		WaitFn:       waitForDahuGit,
	}
	var dahuGit container.ContainerInstance
	var err error
	dahuGit, err = dockerCli.StartContainer(ctx, dahuGitConf)
	if err != nil {
		return http.StatusInternalServerError
	}

	stopOptions := container.ContainerStopOptions{Force: true, RemoveVolumes: true}
	req := buildCheckCloneRequest(gitConfig)
	var reqStatus int
	reqStatus, err = doClone(dahuGit.Ip, req)
	if err != nil {
		// Don't forget to stop the container anyway !
		dockerCli.StopContainer(ctx, dahuGit.Id, stopOptions)
		return http.StatusInternalServerError
	}

	// for the moment, the choice is make
	// to consider that the result of the container
	// stop must be show to the user as a system error.
	// this is because if we have trouble to stop a simple
	// container, then the whole system is in deep trouble.
	// Meaning there is a need for deeper investigations.
	err = dockerCli.StopContainer(ctx, dahuGit.Id, stopOptions)
	if err != nil {
		return http.StatusInternalServerError
	}

	return reqStatus
}

func waitForDahuGit() error {
	// todo create a /status endpoint on dahu-git
	return nil
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

func doClone(ip string, req types.CloneRequest) (int, error) {
	defaultStatus := http.StatusInternalServerError
	body, err := json.Marshal(req)
	if err != nil {
		return defaultStatus, err
	}

	resp, err := http.Post(fmt.Sprintf("http://%s", ip), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return defaultStatus, err
	}
	return resp.StatusCode, nil
}
