package model

import (
	"github.com/jeromedoucet/dahu/core/scm"
)

type HttpAuthConfig struct {
	Url      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"` // todo hide that ! (dont't show when get job)
}

func (a HttpAuthConfig) IsValid() bool {
	if a.Url == "" {
		return false
	}
	return true
}

type SshAuthConfig struct {
	Url         string `json:"url"`
	Key         string `json:"key"`         // todo hide that ! (dont't show when get job)
	KeyPassword string `json:"keyPassword"` // todo hide that ! (dont't show when get job)
}

func (a SshAuthConfig) IsValid() bool {
	if a.Url == "" {
		return false
	} else if a.Key == "" {
		return false
	}
	return true
}

type GitConfig struct {
	HttpAuth *HttpAuthConfig `json:"httpAuth"`
	SshAuth  *SshAuthConfig  `json:"sshAuth"`
}

func (g GitConfig) CheckCredentials() scm.ScmError {
	git := scm.GitInstance
	if g.HttpAuth != nil {
		// todo add some little units test here for rejections cases
		if g.HttpAuth.User == "" || g.HttpAuth.Password == "" {
			return git.CheckConnectionWithoutAuth(g.HttpAuth.Url)
		} else {
			return git.CheckConnectionWithIdAndPassword(g.HttpAuth.Url, g.HttpAuth.User, g.HttpAuth.Password)
		}
	} else if g.SshAuth != nil {
		return git.CheckConnectionWithPrivateKey(g.SshAuth.Url, g.SshAuth.Key, g.SshAuth.KeyPassword)
	} else {
		panic("not implemented yet")
	}
}

// todo test me
func (g GitConfig) IsValid() bool {
	if (g.HttpAuth == nil && g.SshAuth == nil) || (g.HttpAuth != nil && g.SshAuth != nil) {
		return false
	}
	if g.HttpAuth != nil {
		return g.HttpAuth.IsValid()
	} else {
		return g.SshAuth.IsValid()
	}
	return true
}
