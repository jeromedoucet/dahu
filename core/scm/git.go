package scm

// TODO : check uncover test case
import (
	"fmt"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	ssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
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
	fmt.Println(err)
	return fromGitToScmError(err)
}

func (r gitRepository) CheckConnectionWithPrivateKey(url string, key string, keyPassword string) ScmError {
	auth, sshError := ssh.NewPublicKeys("git", []byte(key), keyPassword)
	if sshError != nil {
		return newScmError(sshError.Error(), SshKeyReadingError)
	}
	_, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:        url,
		NoCheckout: true,
		Auth:       auth,
	})
	fmt.Println(err)
	return fromGitToScmError(err)
}

func fromGitToScmError(err error) ScmError {
	if err == nil {
		return nil
	}
	errStr := err.Error()
	switch err {
	case transport.ErrRepositoryNotFound:
		return newScmError(errStr, RepositoryNotFound)
	case transport.ErrAuthenticationRequired:
		return newScmError(errStr, BadCredentials)
	default:
		if strings.Contains(errStr, "no supported methods remain") {
			// this error come directly from clientAuthenticate
			// in client_auth.go from 'golang.org/x/crypto/ssh' package.
			//
			// This error generaly means that the private key used for authentication
			// is not the right one. TODO: maybee use a more specific error ?
			//
			// Because the error is thrown through fmt.Errorf function
			// there is no possibility but checking the text of the error
			// to detect it !
			return newScmError(errStr, BadCredentials)
		} else if strings.Contains(strings.ToLower(errStr), "repository does not exist") {
			// try to catch ssh error related to inexistant repository.
			// This kind of error may be treat as 'unknow error' by the underlying
			// git library (go-git-v4 plumbing/transport/internal/common/common.go)
			//
			// Some Pr may improve this handling, but some case may be missing, that's why
			// we make a try to handle the missing cases here
			return newScmError(errStr, RepositoryNotFound)

		} else {
			return newScmError(errStr, OtherError)
		}
	}
}
