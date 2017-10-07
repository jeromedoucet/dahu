package persistence_test

import (
	"context"
	"testing"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
)

func TestShouldCreateANewJob(t *testing.T) {
	// given
	j := model.Job{Name: "test", Url: "github.com/test"}
	c := configuration.InitConf()

	ctx := context.Background()
	r := persistence.GetRepository(c)
	// todo don't save the database file.

	// when
	nj, err := r.CreateJob(&j, ctx)

	// then
	if err != nil {
		t.Errorf("expect job creation test to have no error but got %s", err.Error())
	}
	if nj.Id == "" {
		t.Errorf("expect CreateJob to affect an Id to the new job but got \"\"")
	}
	if nj.Name != j.Name {
		t.Errorf("expect the new job name to be %s but got %s", j.Name, nj.Name)
	}
	if nj.Url != j.Url {
		t.Errorf("expect the new job url to be %s but got %s", j.Url, nj.Url)
	}
	close(c.Close)
	r.WaitClose()
}
