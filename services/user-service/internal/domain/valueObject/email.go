package valueobject

import "net/mail"

type Email string

func NewEmail(email string) Email {
	return Email(email)
}

func (e Email) String() string {
	return string(e)
}

func (e Email) Validate() error {
	_, err := mail.ParseAddress(string(e))
	return err
}
