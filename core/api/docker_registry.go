package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jeromedoucet/dahu/core/container"
	"github.com/jeromedoucet/dahu/core/model"
)

func (a *Api) handleDockerRegistryCheck(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tokErr := a.checkToken(r)
	if tokErr != nil {
		log.Printf("WARN >> handleGitRepositories encounter error : %s ", tokErr.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
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
