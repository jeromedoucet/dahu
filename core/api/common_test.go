package api_test

import (
	"net/http/httptest"

	"github.com/jeromedoucet/dahu/configuration"
)

// common var used in api_test package
var conf *configuration.Conf
var gitRepoIp string
var s *httptest.Server
