package usecase

import (
	"context"

	"github.com/phongloihong/go-shop/services/user-service/internal/domain/entity"
	"github.com/phongloihong/go-shop/services/user-service/internal/domain/repository"
	"github.com/phongloihong/go-shop/services/user-service/internal/domain/service"
	"github.com/phongloihong/go-shop/services/user-service/internal/usecase/dto"
)

type UserUseCase struct {
	userRepo    repository.UserRepository
	authService service.AuthService
}

func NewUserUseCase(repo repository.UserRepository, authService service.AuthService) *UserUseCase {
	return &UserUseCase{
		userRepo:    repo,
		authService: authService,
	}
}

func (u *UserUseCase) RegisterUser(ctx context.Context, params dto.RegisterRequest) (*entity.User, error) {
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

func (u *UserUseCase) Login(ctx context.Context, params dto.LoginRequest) (*service.TokenPairs, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return nil, err
	}

	if err := user.Password.CompareHash(params.Password); err != nil {
		return nil, err
	}

	ret, err := u.authService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
