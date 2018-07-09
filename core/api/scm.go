package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/scm"
)

func (a *Api) handleGitRepositories(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var gitConfig model.GitConfig
	d := json.NewDecoder(r.Body)
	d.Decode(&gitConfig)
	err := gitConfig.CheckCredentials()
	fmt.Println(err)
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else if err.ErrorType() == scm.RepositoryNotFound {
		w.WriteHeader(http.StatusNotFound)
	} else if err.ErrorType() == scm.BadCredentials {
		w.WriteHeader(http.StatusForbidden)
	} else {
		w.WriteHeader(http.StatusBadRequest) // todo change it when nedeed
	}
}
