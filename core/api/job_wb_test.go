package api

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
)

type mockRepository struct {
	persistence.Repository
	getJobError error
}

func (m *mockRepository) GetJob(id []byte, ctx context.Context) (*model.Job, error) {
	return nil, m.getJobError
}

type mockResponse struct {
	http.ResponseWriter
	code int
}

func (m *mockResponse) WriteHeader(code int) {
	m.code = code
}

type mockRunEngine struct {
}

func (re *mockRunEngine) StartOneRun(job *model.Job, ctx context.Context) error {
	return errors.New("an error")
}

func (r *mockRunEngine) WaitClose() {

}

func TestHandleJobShouldReturn404WhenNoJobFound(t *testing.T) {
	// given
	a := new(Api)
	mock := new(mockRepository)
	mock.getJobError = errors.New("no found")
	a.repository = mock
	url := new(url.URL)
	url.Path = "/jobs/3/trigger"
	req := http.Request{URL: url}
	res := new(mockResponse)

	// when
	a.handleJob(context.Background(), res, &req)

	// then
	if res.code != http.StatusNotFound {
		t.Errorf("expect the http code to be 404 when no Job, but got %d", res.code)
	}
}

func TestHandleJobShouldReturn404WhenJobRunError(t *testing.T) {
	// given
	a := new(Api)
	a.repository = new(mockRepository)
	a.runEngine = new(mockRunEngine)
	url := new(url.URL)
	url.Path = "/jobs/3/trigger"
	req := http.Request{URL: url}
	res := new(mockResponse)

	// when
	a.handleJob(context.Background(), res, &req)

	// then
	if res.code != http.StatusInternalServerError {
		t.Errorf("expect the http code to be 500 when error during job starting, but got %d", res.code)
	}
}
