package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/scm"
)

func (a *Api) handleGitRepositories(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var gitConfig model.GitConfig
	d := json.NewDecoder(r.Body)
	d.Decode(&gitConfig)
	err := gitConfig.CheckCredentials()
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		body := fromErrorToJson(err)
		w.Header().Set("Content-Type", "text/plain")
		if err.ErrorType() == scm.RepositoryNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else if err.ErrorType() == scm.BadCredentials {
			w.WriteHeader(http.StatusForbidden)
		} else if err.ErrorType() == scm.SshKeyReadingError {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write(body)
	}
}
