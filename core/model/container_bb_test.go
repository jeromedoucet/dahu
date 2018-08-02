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
	if registry.User != "" {
		t.Fatalf("Expect ToPublicModel to hide User, but got %s", registry.User)
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
	if registry.Id == nil {
		t.Errorf("expect the Id to have been generated, but is nil")
	}
}

func TestJobIdGenerationFailed(t *testing.T) {
	// given
	id := []byte("existingId")
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
