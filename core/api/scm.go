package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/scm"
)

func (a *Api) handleGitRepositories(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tokErr := a.checkToken(r)
	if tokErr != nil {
		log.Printf("WARN >> handleGitRepositories encounter error : %s ", tokErr.Error())
		w.WriteHeader(http.StatusUnauthorized)
	}
	var gitConfig model.GitConfig
	d := json.NewDecoder(r.Body)
	d.Decode(&gitConfig)
	err := gitConfig.CheckCredentials()
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		if err.ErrorType() == scm.RepositoryNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else if err.ErrorType() == scm.BadCredentials {
			w.WriteHeader(http.StatusForbidden)
		} else if err.ErrorType() == scm.SshKeyReadingError {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
