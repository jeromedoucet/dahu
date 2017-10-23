package model

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// test the encryption of passwords.
func TestSetPasswordShouldUseBcryptWithSalt(t *testing.T) {
	// given
	usr := User{Login: "test"}
	password := []byte("test_pwd_that_must_be_long_enought")

	// when
	usr.SetPassword(password)

	// then
	var err error
	err = bcrypt.CompareHashAndPassword(usr.password, password)
	if err != nil {
		t.Errorf("expect the password to match but get %s", err.Error())
	}
	var cost int
	cost, err = bcrypt.Cost(usr.password)
	if err != nil {
		t.Fatal(err.Error())
	}
	if cost != 13 {
		t.Errorf("expect the complexity of the hash to be 13 but got %d", cost)
	}
}

// test the encryption of passwords when no
// password passed to the function.
func TestSetPasswordShouldFailedWhenNoPassword(t *testing.T) {
	// given
	usr := User{Login: "test"}
	password := []byte("")

	// when
	err := usr.SetPassword(password)

	// then
	if err == nil {
		t.Error("expect an error when no password but got nil")
	}
}

func TestComparePasswordFailure(t *testing.T) {
	// given
	password := []byte("some_password")
	hashedPassword, _ := bcrypt.GenerateFromPassword(password, 13)
	usr := User{Login: "test", password: hashedPassword}

	// when
	err := usr.ComparePassword([]byte("other_password"))

	if err == nil {
		t.Error("expect an error when trying to compare a wrong password")
	}
}

func TestComparePasswordSuccess(t *testing.T) {
	// given
	password := []byte("some_password")
	hashedPassword, _ := bcrypt.GenerateFromPassword(password, 13)
	usr := User{Login: "test", password: hashedPassword}

	// when
	err := usr.ComparePassword(password)

	if err != nil {
		t.Errorf("expect no error when trying to compare a right password, but got %+v", err)
	}
}
