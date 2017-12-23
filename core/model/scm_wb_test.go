package model

import (
	"testing"
)

func TestGetImageShouldReturnGitImage(t *testing.T) {
	// given
	scm := Scm{Type: GIT}

	// when
	image := scm.getImage()

	// then
	if image != "dahuci/git" {
		t.Errorf("expect to have alpine/git but got %s", image)
	}
}

func TestGetImageShouldReturnSvnImage(t *testing.T) {
	// given
	scm := Scm{Type: SVN}

	// when
	image := scm.getImage()

	// then
	if image != "dahuci/svn" {
		t.Errorf("expect to have alpine/git but got %s", image)
	}
}
