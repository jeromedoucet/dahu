package model

import (
	"github.com/jeromedoucet/dahu/core/scm"
)

type ScmType int

const (
	GIT ScmType = 1 + iota
	SVN
)

// todo remove that type when refactoring the job
type Scm struct {
	RepoUrl string
	Type    ScmType
}

func (s Scm) getImage() string {
	switch s.Type {
	case GIT:
		return "dahuci/git"
	case SVN:
		return "dahuci/svn"
	default:
		return ""
	}
}

type HttpAuthConfig struct {
	Url      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type SshAuthConfig struct {
	Url string
}

type GitConfig struct {
	HttpAuth *HttpAuthConfig `json:"httpAuth"`
	SshAuth  *SshAuthConfig  `json:"sshAuth"`
}

func (g GitConfig) CheckCredentials() scm.ScmError {
	git := scm.GitInstance
	if g.HttpAuth != nil {
		// todo think to add some little units test here for rejections cases
		// todo fix demeter law violation
		if g.HttpAuth.User == "" || g.HttpAuth.Password == "" {
			return git.CheckConnectionWithoutAuth(g.HttpAuth.Url)
		} else {
			return git.CheckConnectionWithIdAndPassword(g.HttpAuth.Url, g.HttpAuth.User, g.HttpAuth.Password)
		}
	} else {
		return nil
	}
}
