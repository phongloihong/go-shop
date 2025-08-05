package valueobject

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Password string

func NewPassword(password string) Password {
	return Password(password)
}

func (p Password) String() string {
	return string(p)
}

func (p Password) Validate() error {
	if len(p) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func (p Password) Hash() (string, error) {
	bytes, error := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(bytes), error
}

func (p Password) CompareHash(passwordString string) error {
	err := bcrypt.CompareHashAndPassword([]byte(p), []byte(passwordString))
	if err == nil {
		return nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return fmt.Errorf("password does not match: %w", err)
	}

	return err
}
