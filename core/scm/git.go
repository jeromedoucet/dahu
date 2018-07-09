package scm

import (
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type gitRepository struct {
}

func (r gitRepository) CheckConnectionWithoutAuth(url string) ScmError {
	_, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:        url,
		NoCheckout: true,
	})
	return fromGitToScmError(err)
}

func (r gitRepository) CheckConnectionWithIdAndPassword(url string, id string, password string) ScmError {
	auth := &http.BasicAuth{Username: id, Password: password}
	_, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:        url,
		NoCheckout: true,
		Auth:       auth,
	})
	return fromGitToScmError(err)
}

func fromGitToScmError(err error) ScmError {
	if err == nil {
		return nil
	}
	switch err {
	case transport.ErrRepositoryNotFound:
		return newScmError(err.Error(), RepositoryNotFound)
	case transport.ErrAuthenticationRequired:
		return newScmError(err.Error(), BadCredentials)
	default:
		return newScmError(err.Error(), OtherError)
	}
}
