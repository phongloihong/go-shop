package connect

import (
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/phongloihong/go-shop/services/user-service/external/gen/user/v1/userv1connect"
	"github.com/phongloihong/go-shop/services/user-service/internal/config"
	"github.com/phongloihong/go-shop/services/user-service/internal/infrastructure/auth"
	"github.com/phongloihong/go-shop/services/user-service/internal/infrastructure/database/postgres"
	"github.com/phongloihong/go-shop/services/user-service/internal/infrastructure/database/postgres/sqlc"
	"github.com/phongloihong/go-shop/services/user-service/internal/usecase"
)

func StartConnect(cfg *config.Config, dbConn sqlc.DBTX) *http.Server {
	mux := http.NewServeMux()

	// create interceptors
	interceptors := connect.WithInterceptors(
		newRecoverInterceptors(),
	)

	authService := auth.NewJWTService(
		[]byte(cfg.Auth.AccessSecret),
		[]byte(cfg.Auth.RefreshSecret),
		time.Duration(30*time.Minute), // expires in 30 minutes
		time.Duration(7*24*time.Hour), // expires in 7 days
	)

	userRepo := postgres.NewUserRepository(dbConn)
	userUseCase := usecase.NewUserUseCase(userRepo, authService)
	userHandler := NewUserServiceHandler(userUseCase)
	mux.Handle(userv1connect.NewUserServiceHandler(userHandler, interceptors))

	return &http.Server{Handler: mux}
}
