package repository

import "travel-planning/internal/domain/entity"

type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	GetUserByID(id string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetPublicProfileByIds(ids []string) ([]*entity.UserPublicProfile, error)
}
