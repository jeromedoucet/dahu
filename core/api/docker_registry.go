package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jeromedoucet/dahu/core/container"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
	"github.com/jeromedoucet/route"
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

func (a *Api) handleDockerRegistry(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		a.onDockerRegistryGet(ctx, w, r)
	} else if r.Method == http.MethodDelete {
		a.onDockerRegistryDelete(ctx, w, r)
	} else {
		// todo return appropriate http code with a corresponding test
	}
}

// http handler that deals with get request on a single docker registry resource
func (a *Api) onDockerRegistryGet(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	path := route.SplitPath(r.URL.Path)
	registryId := path[len(path)-1]
	registry, persistenceErr := a.repository.GetDockerRegistry([]byte(registryId), ctx)
	if persistenceErr != nil {
		log.Printf("ERROR >> onDockerRegistryGet encounter error : %s", persistenceErr.Error())
		body := fromErrorToJson(persistenceErr)
		if persistenceErr.ErrorType() == persistence.NotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(body)
		return
	}
	registry.ToPublicModel()
	body, err := json.Marshal(registry)
	if err != nil {
		log.Printf("ERROR >> onDockerRegistryGet encounter error : %s", err.Error())
		body := fromErrorToJson(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(body)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// http handler that deals with delete request on a docker registry resource
func (a *Api) onDockerRegistryDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	path := route.SplitPath(r.URL.Path)
	registryId := path[len(path)-1]
	persistenceErr := a.repository.DeleteDockerRegistry([]byte(registryId))
	if persistenceErr != nil {
		log.Printf("ERROR >> onDockerRegistryDelete encounter error : %s", persistenceErr.Error())
		body := fromErrorToJson(persistenceErr)
		if persistenceErr.ErrorType() == persistence.NotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(body)
		return
	}
	w.WriteHeader(http.StatusOK)
}
