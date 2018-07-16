package api

import (
	"net/http"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/persistence"
	"github.com/jeromedoucet/dahu/core/run"
	"github.com/jeromedoucet/route"
)

type Api struct {
	conf       *configuration.Conf
	router     *route.DynamicRouter
	repository persistence.Repository
	runEngine  run.RunEngine
}

func (a *Api) initRouter() {
	a.router = route.NewDynamicRouter()
	a.router.HandleFunc("/jobs", a.handleJobs)
	a.router.HandleFunc("/jobs/:jobId/run", a.handleJob)
	a.router.HandleFunc("/login", a.handleAuthentication)
	a.router.HandleFunc("/scm/git/repository", a.handleGitRepositories)
}

func (a *Api) Handler() http.Handler {
	return a.router
}

// todo pass a context for timeout
func (a *Api) Close() {
	// close the run engine
	// then wait for the repository
	// to greceful shutdown
	a.runEngine.WaitClose()
	a.repository.WaitClose()
}

func InitRoute(c *configuration.Conf) *Api {
	a := new(Api)
	a.conf = c
	a.repository = persistence.GetRepository(c)
	a.initRouter()
	a.runEngine = run.NewRunEngine(c)
	return a
}
