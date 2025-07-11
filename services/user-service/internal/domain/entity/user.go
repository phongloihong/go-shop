package entity

import (
	valueobject "travel-planning/internal/domain/valueObject"
	"travel-planning/internal/pkg/utils"
)

type User struct {
	ID        string                `json:"id"`
	FirstName string                `json:"first_name"`
	LastName  string                `json:"last_name"`
	Email     valueobject.Email     `json:"email"`
	Phone     valueobject.Phone     `json:"phone"`
	Password  valueobject.Passwword `json:"-"`
	CreatedAt valueobject.DateTime  `json:"created_at"`
	UpdatedAt valueobject.DateTime  `json:"updated_at"`
}

func NewUser(firstName, lastName, email, phone, password, createdAt, updatedAt string) (*User, error) {
	passwordVO := valueobject.NewPassword(password)
	emailVO := valueobject.NewEmail(email)
	phoneVO := valueobject.NewPhone(phone)
	nowVO := valueobject.NewTime(utils.TimeNow())

	user := &User{
		ID:        utils.NewUUID(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     emailVO,
		Phone:     phoneVO,
		Password:  passwordVO,
		CreatedAt: nowVO,
		UpdatedAt: nowVO,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) Validate() error {
	if err := u.Email.Validate(); err != nil {
		return err
	}

	if err := u.Password.Validate(); err != nil {
		return err
	}

	if err := u.Phone.Validate(); err != nil {
		return err
	}

	return nil
}
