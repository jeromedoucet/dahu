package container

import (
	"context"
	"errors"
	"testing"

	client "github.com/docker/docker/client"
)

func TestCheckRegistryConnectionClientIssue(t *testing.T) {
	// given
	errorMsg := "one error sir !"
	docker := dockerClient{clientOpts: func(cli *client.Client) error { return errors.New(errorMsg) }}
	ctx := context.Background()

	// when
	err := docker.CheckRegistryConnection(ctx, RegistryBasicConf{})

	// then
	if err == nil {
		t.Fatal("expect having an error when trying to CheckRegistryConnection with wrong api version but got nil")
	}
	if err.ErrorType() != OtherError {
		t.Fatalf("expect having a ContainerError with type %d but got %d", OtherError, err.ErrorType())
	}
	if err.Error() != errorMsg {
		t.Fatalf("expect having a ContainerError with msg %s but got %s", errorMsg, err.Error())
	}
}
