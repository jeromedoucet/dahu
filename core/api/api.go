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
	a.router.HandleFunc("/jobs", a.handleJobs)
}

func (a *Api) Handler() http.Handler {
	return a.router
}

func InitRoute(c *configuration.Conf) *Api {
	a := new(Api)
	a.conf = c
	a.repository = persistence.GetRepository(c)
	a.initRouter()
	return a
}
