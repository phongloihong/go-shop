package usecase

import (
	"context"

	"github.com/phongloihong/go-shop/services/user-service/internal/domain/entity"
	"github.com/phongloihong/go-shop/services/user-service/internal/domain/repository"
	"github.com/phongloihong/go-shop/services/user-service/internal/usecase/dto"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: repo,
	}
}

func (u *UserUseCase) RegisterUser(ctx context.Context, params dto.RegisterUserRequest) (*entity.User, error) {
	// Create entity
	newUser, err := entity.NewUser(
		params.FirstName,
		params.LastName,
		params.Email,
		params.Phone,
		params.Password,
	)
	if err != nil {
		return nil, err
	}

	// save to database
	ret, err := u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
