package persistence_test

import (
	"context"
	"testing"

	"github.com/jeromedoucet/dahu/configuration"
	"github.com/jeromedoucet/dahu/core/model"
	"github.com/jeromedoucet/dahu/core/persistence"
	"github.com/jeromedoucet/dahu/tests"
)

// todo test unicity for user

// test to ensure that a default user is inserted
// when create the db
func TestInsertDefaultUser(t *testing.T) {
	// given
	c := configuration.InitConf()
	ctx := context.Background()
	login := "dahu"

	// when
	r := persistence.GetRepository(c)
	actualUser, err := r.GetUser(login, ctx)

	// close and remove the db
	tests.CleanPersistence(c)

	// then
	if err != nil {
		t.Errorf("expect to have no error when finding existing user, but got %s", err.Error())
	}
	if actualUser.Login != login {
		t.Errorf("expect to get user %s but got %s", login, actualUser.String())
	}
}

// test the nominal case of #GetUser
func TestGetUserShouldReturnTheUserWhenItExists(t *testing.T) {
	// given
	u := model.User{Login: "test"}
	u.SetPassword([]byte("test_test_test_test"))
	c := configuration.InitConf()

	ctx := context.Background()
	tests.InsertObject(c, []byte("users"), []byte(u.Login), u)
	rep := persistence.GetRepository(c)

	// when
	actualUser, err := rep.GetUser(u.Login, ctx)

	// close and remove the db
	tests.CleanPersistence(c)

	// then
	if err != nil {
		t.Errorf("expect to have no error when finding existing user, but got %s", err.Error())
	}
	if actualUser.Login != u.Login {
		t.Errorf("expect to get user %s but got %s", u.String(), actualUser.String())
	}
}

// test the case when the user is not found for #GetUser.
// => an error is returned
func TestGetUserShouldReturnAnErrorWhenItDoesntExist(t *testing.T) {
	// given
	u := model.User{Login: "test"}
	u.SetPassword([]byte("test_test_test_test"))
	c := configuration.InitConf()

	ctx := context.Background()
	rep := persistence.GetRepository(c)

	// when
	actualUser, err := rep.GetUser(u.Login, ctx)

	// close and remove the db
	tests.CleanPersistence(c)

	// then
	if err == nil {
		t.Error("expect to have an error when searching non-existing user, but got nil")
	}
	if actualUser != nil {
		t.Errorf("expect to get nil but got %s", actualUser.String())
	}
}

// test that the case where there is no bucket for user
// is properly handle at #GetUser
func TestGetUserShouldReturnAnErrorWhenNoBucket(t *testing.T) {
	// todo move it in a 'white box' to make it possible (with no deadlock)
	t.SkipNow()
	// given
	u := model.User{Login: "test"}
	u.SetPassword([]byte("test_test_test_test"))
	c := configuration.InitConf()
	tests.DeleteBucket(c, []byte("users"))

	ctx := context.Background()
	rep := persistence.GetRepository(c)

	// when
	actualUser, err := rep.GetUser(u.Login, ctx)

	// close and remove the db
	tests.CleanPersistence(c)

	// then
	if err == nil {
		t.Error("expect to have an error when searching user without any users bucket, but got nil")
	}
	if actualUser != nil {
		t.Errorf("expect to get nil but got %s", actualUser.String())
	}
}
