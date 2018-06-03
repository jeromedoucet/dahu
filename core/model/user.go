package model

import (
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var regexPassword *regexp.Regexp = regexp.MustCompile(".{12}")

// type used for authentication
// on /login
type Login struct {
	Id       string `json:"id"`
	Password []byte `json:"password"`
}

// this is the answer to
// a successfull login call
type Token struct {
	Value string `json:"value"`
}

// User of Dahu system.
// A user can be a human, or not.
// It represents only an identity.
type User struct {
	Login    string `json:"login"`
	Password []byte // unexported and excluded from json marshaling
}

func (u *User) String() string {
	return fmt.Sprintf("{Login:%s}", u.Login)
}

// hash and affect a new password
func (u *User) SetPassword(password []byte) error {
	if !regexPassword.Match(password) {
		return errors.New("expect a password with at least 12 characters")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 13)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

func (u *User) ComparePassword(password []byte) error {
	return bcrypt.CompareHashAndPassword(u.Password, password)
}
