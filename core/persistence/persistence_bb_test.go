package persistence_test

import (
	"context"
	"testing"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
	"github.com/jeromedoucet/dahu/tests"
)

// Nominal test of job creation for
// inmemory db
func TestShouldCreateANewJobWithInMemoryDb(t *testing.T) {
	// given
	j := model.Job{Name: "test", Url: "github.com/test"}
	c := configuration.InitConf()

	ctx := context.Background()
	r := persistence.GetRepository(c)
	// todo don't save the database file.

	// when
	nj, err := r.CreateJob(&j, ctx)

	// close and remove the db
	tests.CleanPersistence(c)

	// then
	if err != nil {
		t.Fatalf("expect job creation test to have no error but got %s", err.Error())
	}
	if string(nj.Id) == "" {
		t.Errorf("expect CreateJob to affect an Id to the new job but got \"\"")
	}
	if nj.Name != j.Name {
		t.Errorf("expect the new job name to be %s but got %s", j.Name, nj.Name)
	}
	if nj.Url != j.Url {
		t.Errorf("expect the new job url to be %s but got %s", j.Url, nj.Url)
	}
}
