package valueobject

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Passwword string

func NewPassword(password string) Passwword {
	return Passwword(password)
}

func (p Passwword) String() string {
	return string(p)
}

func (p Passwword) Validate() error {
	if len(p) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func (p Passwword) Hash() (string, error) {
	bytes, error := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(bytes), error
}

func (p Passwword) CompareHash(hash Passwword) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
}
