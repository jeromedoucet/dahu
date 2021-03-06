package model_test

import (
	"testing"

	"github.com/jeromedoucet/dahu/core/model"
)

func TestIsValidNoAuth(t *testing.T) {
	// given
	conf := model.GitConfig{}

	// when
	isValid := conf.IsValid()

	// then
	if isValid {
		t.Error("expect GitConfig to be invalid when no auth conf, but is valid")
	}
}

func TestIsValidHttpAuthNoUrl(t *testing.T) {
	// given
	httpAuth := model.HttpAuthConfig{User: "somebody", Password: "some-password"}
	conf := model.GitConfig{HttpAuth: &httpAuth}

	// when
	isValid := conf.IsValid()

	// then
	if isValid {
		t.Error("expect GitConfig to be invalid when http auth without url, but is valid")
	}
}

func TestIsValidHttpAuth(t *testing.T) {
	// given
	httpAuth := model.HttpAuthConfig{Url: "http://some-domain/some-repo"}
	conf := model.GitConfig{HttpAuth: &httpAuth}

	// when
	isValid := conf.IsValid()

	// then
	if !isValid {
		t.Error("expect GitConfig to be valid when correct http auth conf, but is invalid")
	}
}

func TestIsValidHttpAuthWithUserAndPwd(t *testing.T) {
	// given
	httpAuth := model.HttpAuthConfig{Url: "http://some-domain/some-repo", User: "user", Password: "password"}
	conf := model.GitConfig{HttpAuth: &httpAuth}

	// when
	isValid := conf.IsValid()

	// then
	if !isValid {
		t.Error("expect GitConfig to be valid when correct http auth conf, but is invalid")
	}
}

func TestIsValidSshAuthNoUrl(t *testing.T) {
	// given
	sshAuth := model.SshAuthConfig{Key: "some-private-key", KeyPassword: "some-password"}
	conf := model.GitConfig{SshAuth: &sshAuth}

	// when
	isValid := conf.IsValid()

	// then
	if isValid {
		t.Error("expect GitConfig to be invalid when ssh auth without url, but is valid")
	}
}

func TestIsValidSshAuthNoKey(t *testing.T) {
	// given
	sshAuth := model.SshAuthConfig{Url: "git@some-domain/some-repo.git", KeyPassword: "some-password"}
	conf := model.GitConfig{SshAuth: &sshAuth}

	// when
	isValid := conf.IsValid()

	// then
	if isValid {
		t.Error("expect GitConfig to be invalid when ssh auth without key, but is valid")
	}
}

func TestIsValidSshAuthWithPwd(t *testing.T) {
	// given
	sshAuth := model.SshAuthConfig{Url: "git@some-domain/some-repo.git", Key: "some-key", KeyPassword: "some-password"}
	conf := model.GitConfig{SshAuth: &sshAuth}

	// when
	isValid := conf.IsValid()

	// then
	if !isValid {
		t.Error("expect GitConfig to be valid when ssh auth with password, but is valid")
	}
}

func TestIsValidSshAuthWithoutPwd(t *testing.T) {
	// given
	sshAuth := model.SshAuthConfig{Url: "git@some-domain/some-repo.git", Key: "some-key"}
	conf := model.GitConfig{SshAuth: &sshAuth}

	// when
	isValid := conf.IsValid()

	// then
	if !isValid {
		t.Error("expect GitConfig to be valid when ssh auth without password, but is valid")
	}
}

func TestIsValidWithTwoAuth(t *testing.T) {
	// given
	sshAuth := model.SshAuthConfig{Url: "git@some-domain/some-repo.git", Key: "some-key"}
	httpAuth := model.HttpAuthConfig{Url: "http://some-domain/some-repo"}
	conf := model.GitConfig{SshAuth: &sshAuth, HttpAuth: &httpAuth}

	// when
	isValid := conf.IsValid()

	// then
	if isValid {
		t.Error("expect GitConfig to be invalid when ssh auth and http auth both configured, but is valid")
	}
}

func TestToPublicModelSshAuthConfig(t *testing.T) {
	// given
	sshAuth := model.SshAuthConfig{Url: "git@some-domain/some-repo.git", Key: "some-key", KeyPassword: "some key password"}

	// when
	sshAuth.ToPublicModel()

	// then
	if sshAuth.Url != "git@some-domain/some-repo.git" {
		t.Error("expect ToPublicModel to let Url untouched, but it has changed")
	}

	if sshAuth.Key != "" {
		t.Error("expect ToPublicModel to clean Key, but it is not empty")
	}

	if sshAuth.KeyPassword != "" {
		t.Error("expect ToPublicModel to clean KeyPassword, but it is not empty")
	}
}

func TestToPublicModelHttpAuthConfig(t *testing.T) {
	// given
	httpAuth := model.HttpAuthConfig{Url: "http://some-domain/some-repo", User: "some user", Password: "some password"}

	// when
	httpAuth.ToPublicModel()

	// then
	if httpAuth.Url != "http://some-domain/some-repo" {
		t.Error("expect ToPublicModel to let Url untouched, but it has changed")
	}

	if httpAuth.User != "" {
		t.Error("expect ToPublicModel to clean User, but it is not empty")
	}

	if httpAuth.Password != "" {
		t.Error("expect ToPublicModel to clean Password, but it is not empty")
	}
}
