package model_test

import (
	"errors"
	"testing"

	"github.com/jeromedoucet/dahu/core/model"
)

func TestIsOutDatedReturnTrue(t *testing.T) {
	// given
	err := model.NewOutDated("for test")

	// when
	res := model.IsOutDated(err)

	// then
	if !res {
		t.Error("expect #IsOutDated to return true but got false")
	}
}

func TestIsOutDatedReturnFalse(t *testing.T) {
	// given
	err := errors.New("for test")

	// when
	res := model.IsOutDated(err)

	// then
	if res {
		t.Error("expect #IsOutDated to return fqlse but got true")
	}
}

func TestIsNoMorePersistedReturnTrue(t *testing.T) {
	// given
	err := model.NewNoMorePersisted("for test")

	// when
	res := model.IsNoMorePersisted(err)

	// then
	if !res {
		t.Error("expect #IsNoMorePersisted to return true but got false")
	}
}

func TestIsNoMorePersistedReturnFalse(t *testing.T) {
	// given
	err := errors.New("for test")

	// when
	res := model.IsNoMorePersisted(err)

	// then
	if res {
		t.Error("expect #IsNoMorePersisted to return false but got true")
	}
}
