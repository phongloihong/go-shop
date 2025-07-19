package repository

import (
	"context"

	"github.com/phongloihong/go-shop/services/user-service/internal/domain/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (int64, error)
	ChangePassword(ctx context.Context, id string, newPassword string) (int64, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetPublicProfileByIds(ctx context.Context, ids []string) ([]*entity.UserPublicProfile, error)
}
