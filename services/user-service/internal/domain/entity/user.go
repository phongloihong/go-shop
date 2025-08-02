package entity

import (
	valueobject "github.com/phongloihong/go-shop/services/user-service/internal/domain/valueObject"
	"github.com/phongloihong/go-shop/services/user-service/internal/pkg/utils"
)

type User struct {
	ID        string               `json:"id"`
	FirstName string               `json:"first_name"`
	LastName  string               `json:"last_name"`
	Email     valueobject.Email    `json:"email"`
	Phone     valueobject.Phone    `json:"phone"`
	Password  valueobject.Password `json:"-"`
	CreatedAt valueobject.DateTime `json:"created_at"`
	UpdatedAt valueobject.DateTime `json:"updated_at"`
}

func NewUser(firstName, lastName, email, phone, password string) (*User, error) {
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

func UserFromDatabase(id, firstName, lastName, email, phone, password string, createdAt, updatedAt int64) *User {
	passwordVO := valueobject.NewPassword(password)
	emailVO := valueobject.NewEmail(email)
	phoneVO := valueobject.NewPhone(phone)
	createdAtVO := valueobject.NewTime(createdAt)
	updatedAtVO := valueobject.NewTime(updatedAt)

	user := &User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     emailVO,
		Phone:     phoneVO,
		Password:  passwordVO,
		CreatedAt: createdAtVO,
		UpdatedAt: updatedAtVO,
	}

	return user
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
