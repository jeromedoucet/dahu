package persistence_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
	"github.com/jeromedoucet/dahu/tests"
)

// test that we may not try to insert / create
// a docker registry that already has an id
func TestCreateDockerRegistryIdExist(t *testing.T) {
	// given
	registry := &model.DockerRegistry{Name: "test", Url: "localhost:5000", User: "tester", Password: "test"}
	registry.GenerateId()
	expectedErrorMsg := fmt.Sprintf("the id %+v already defined", string(registry.Id))
	c := configuration.InitConf()

	ctx := context.Background()
	r := persistence.GetRepository(c)

	// when
	nr, err := r.CreateDockerRegistry(registry, ctx)

	// close and remove the db
	tests.CleanPersistence(c)

	// then
	if nr != nil {
		t.Fatalf(`expect to get no new job for a call on #CreateJob
		with a job that already have an id but got %+v`, nr)
	}
	if err == nil {
		t.Fatal(`expect to have an error when calling #CreateJob with
		a job that already have an id, but got nil`)
	}
	if err.Error() != expectedErrorMsg {
		t.Fatalf("wrong error messager got %s, expected %s", err.Error(), expectedErrorMsg)
	}
}
