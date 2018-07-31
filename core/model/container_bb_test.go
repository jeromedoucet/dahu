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
