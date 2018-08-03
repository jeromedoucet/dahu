package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jeromedoucet/dahu/core/container"
	"github.com/jeromedoucet/dahu/core/model"
)

// Allow to test one docker registry configuration
func (a *Api) handleDockerRegistryCheck(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var registry model.DockerRegistry
	d := json.NewDecoder(r.Body)
	d.Decode(&registry)
	err := registry.CheckCredentials(ctx)
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		body := fromErrorToJson(err)
		if err.ErrorType() == container.RegistryNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else if err.ErrorType() == container.BadCredentials {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write(body)
	}
}

// create a new docker registry. Will fail if there
// is already an id in the given registry
func (a *Api) handleDockerRegistryCreation(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var registry model.DockerRegistry
	var err error
	d := json.NewDecoder(r.Body)
	d.Decode(&registry)
	var newRegistry *model.DockerRegistry
	newRegistry, err = a.repository.CreateDockerRegistry(&registry, ctx)
	if err != nil {
		log.Printf("ERROR >> dockerRegistryCreation encounter error : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newRegistry.ToPublicModel()
	body, _ := json.Marshal(newRegistry)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}
