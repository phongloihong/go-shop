package connect

import (
	"context"

	"connectrpc.com/connect"
	userv1 "github.com/phongloihong/go-shop/services/user-service/external/gen/user/v1"
	domain_error "github.com/phongloihong/go-shop/services/user-service/internal/domain/domain_errors"
	"github.com/phongloihong/go-shop/services/user-service/internal/usecase"
	"github.com/phongloihong/go-shop/services/user-service/internal/usecase/dto"
)

type userServiceHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserServiceHandler(
	userUseCase *usecase.UserUseCase,
) *userServiceHandler {
	return &userServiceHandler{
		userUseCase: userUseCase,
	}
}

func (h *userServiceHandler) Register(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error) {
	params := dto.RegisterRequest{
		FirstName: req.Msg.FirstName,
		LastName:  req.Msg.LastName,
		Email:     req.Msg.Email,
		Phone:     req.Msg.Phone,
		Password:  req.Msg.Password,
	}

	_, err := h.userUseCase.RegisterUser(ctx, params)
	if err != nil {
		return nil, domain_error.MapError(err)
	}

	ret := &userv1.RegisterResponse{
		Success: true,
	}

	return connect.NewResponse(ret), nil
}

func (h *userServiceHandler) Login(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error) {
	ret, err := h.userUseCase.Login(ctx, dto.LoginRequest{
		Email:    req.Msg.Email,
		Password: req.Msg.Password,
	})
	if err != nil {
		return nil, domain_error.MapError(err)
	}

	return connect.NewResponse(&userv1.LoginResponse{
		AccessToken:  ret.AccessToken,
		RefreshToken: ret.RefreshToken,
		ExpiresIn:    ret.ExpiresIn,
	}), nil
}

func (h *userServiceHandler) ChangePassword(ctx context.Context, req *connect.Request[userv1.ChangePasswordRequest]) (*connect.Response[userv1.ChangePasswordResponse], error) {
	return nil, nil
}

func (h *userServiceHandler) GetProfile(ctx context.Context, req *connect.Request[userv1.GetProfileRequest]) (*connect.Response[userv1.GetProfileResponse], error) {
	return nil, nil
}

func (h *userServiceHandler) GetPublicProfile(ctx context.Context, req *connect.Request[userv1.GetPublicProfileRequest]) (*connect.Response[userv1.GetPublicProfileResponse], error) {
	return nil, nil
}
