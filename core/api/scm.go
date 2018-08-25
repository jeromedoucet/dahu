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
	statusCode := scm.CheckClone(ctx, gitConfig)
	w.WriteHeader(statusCode)
}
