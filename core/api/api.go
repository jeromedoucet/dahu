package api

import (
	"net/http"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/persistence"
	"github.com/jeromedoucet/route"
)

type Api struct {
	conf       *configuration.Conf
	router     *route.DynamicRouter
	repository persistence.Repository
}

func (a *Api) initRouter() {
	a.router = route.NewDynamicRouter()
	a.router.HandleFunc("/jobs", a.handleJobs, a.authFilter)
	a.router.HandleFunc("/login", a.handleAuthentication)
	a.router.HandleFunc("/scm/git/repository", a.handleGitRepositories, a.authFilter)
	a.router.HandleFunc("/containers/docker/registries/test", a.handleDockerRegistryCheck, a.authFilter)
	a.router.HandleFunc("/containers/docker/registries", a.handleDockerRegistryCreation, a.authFilter)
}

func (a *Api) Handler() http.Handler {
	return a.router
}

// todo pass a context for timeout
func (a *Api) Close() {
	// wait for the repository
	// to graceful shutdown
	close(a.conf.Close)
	a.repository.WaitClose()
}

func InitRoute(c *configuration.Conf) *Api {
	a := new(Api)
	a.conf = c
	a.repository = persistence.GetRepository(c)
	a.initRouter()
	return a
}
