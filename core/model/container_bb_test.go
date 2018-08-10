package model_test

import (
	"testing"

	"github.com/jeromedoucet/dahu/core/model"
)

func TestRegistryToPublicModel(t *testing.T) {
	// given
	registry := &model.DockerRegistry{Name: "test", Url: "localhost:5000", User: "tester", Password: "test"}

	// when
	registry.ToPublicModel()

	// then
	if registry.Name != "test" {
		t.Fatalf("Expect ToPublicModel to preserve Name, but got %s", registry.Name)
	}
	if registry.Url != "localhost:5000" {
		t.Fatalf("Expect ToPublicModel to preserve Url, but got %s", registry.Url)
	}
	if registry.User != "tester" {
		t.Fatalf("Expect ToPublicModel to preserve User, but got %s", registry.User)
	}
	if registry.Password != "" {
		t.Fatalf("Expect ToPublicModel to hide Password, but got %s", registry.Password)
	}
}

func TestJobIdGenerationSuccessFull(t *testing.T) {
	// given
	registry := &model.DockerRegistry{Name: "test", Url: "localhost:5000", User: "tester", Password: "test"}

	// when
	err := registry.GenerateId()

	// then
	if err != nil {
		t.Errorf("Expect #GenerateId to return nil, but got %v", err)
	}
	if registry.Id == "" {
		t.Errorf("expect the Id to have been generated, but is nil")
	}
}

func TestJobIdGenerationFailed(t *testing.T) {
	// given
	id := "existingId"
	registry := &model.DockerRegistry{Id: id, Name: "test", Url: "localhost:5000", User: "tester", Password: "test"}

	// when
	err := registry.GenerateId()

	// then
	if err == nil {
		t.Errorf("Expect #GenerateId to return an error, but got nil")
	}
	if string(registry.Id) != string(id) {
		t.Errorf("expect the Id not to have changed, but got %s", string(registry.Id))
	}
}

func TestMergeForUpdate(t *testing.T) {
	// given
	registry := &model.DockerRegistry{
		Name:     "name",
		Url:      "https://some-domain/path",
		User:     "some user",
		Password: "some password",
	}
	registry.NewLastModificationTime()
	registryUpdate := new(model.DockerRegistryUpdate)
	registryUpdate.Name = "updated name"
	registryUpdate.Url = "https://some-new-domain/path"
	registryUpdate.User = "some new user"
	registryUpdate.Password = "some password"
	registryUpdate.NewLastModificationTime()
	registryUpdate.ChangedFields = []string{"name", "url", "user", "badField"}

	// when
	mergedRegistry := registryUpdate.MergeForUpdate(registry)

	// then
	if mergedRegistry.Name != registryUpdate.Name {
		t.Fatal("expect the name to have been merged as updated field")
	}
	if mergedRegistry.Url != registryUpdate.Url {
		t.Fatal("expect the url to have been merged as updated field")
	}
	if mergedRegistry.User != registryUpdate.User {
		t.Fatal("expect the user to have been merged as updated field")
	}
	if mergedRegistry.Password != registry.Password {
		t.Fatal("expect the password to have been merged as a non updated field")
	}
	if mergedRegistry.LastModificationTime != registryUpdate.LastModificationTime {
		t.Fatal("expect the LastModificationTime to be the one from registryUpdate")
	}
}
