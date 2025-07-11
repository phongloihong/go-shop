package valueobject

import "fmt"

type Phone string

func NewPhone(phone string) Phone {
	return Phone(phone)
}

func (p Phone) String() string {
	return string(p)
}

func (p Phone) Validate() error {
	if p == "" {
		return nil
	}

	if len(p) < 10 || len(p) > 15 {
		return fmt.Errorf("phone number must be between 10 and 15 characters long")
	}

	return nil
}
