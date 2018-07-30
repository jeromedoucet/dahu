package api_test

import (
	"os"
	"testing"

	"github.com/jeromedoucet/dahu/tests"
)

func TestMain(m *testing.M) {
	gogsId := tests.StartGogs()
	registryId := tests.StartDockerRegistry()
	retCode := m.Run()
	tests.StopContainer(gogsId)
	tests.StopContainer(registryId)
	os.Exit(retCode)
}
