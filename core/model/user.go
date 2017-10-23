package model

import (
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var regexPassword *regexp.Regexp = regexp.MustCompile(".{12}")

// User of Dahu system.
// A user can be a human, or not.
// It represents only an identity.
type User struct {
	Login    string `json:"login"`
	password []byte // unexported and excluded from json marshaling
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
	u.password = hashedPassword
	return nil
}

func (u *User) ComparePassword(password []byte) error {
	return bcrypt.CompareHashAndPassword(u.password, password)
}
