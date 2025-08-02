package connect

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/phongloihong/go-shop/services/user-service/external/gen/user/v1/userv1connect"
	"github.com/phongloihong/go-shop/services/user-service/internal/infrastructure/database/postgres"
	"github.com/phongloihong/go-shop/services/user-service/internal/infrastructure/database/postgres/sqlc"
	"github.com/phongloihong/go-shop/services/user-service/internal/usecase"
)

func StartConnect(dbConn sqlc.DBTX) *http.Server {
	mux := http.NewServeMux()

	// create interceptors
	interceptors := connect.WithInterceptors(
		newRecoverInterceptors(),
	)

	userRepo := postgres.NewUserRepository(dbConn)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := NewUserServiceHandler(userUseCase)
	mux.Handle(userv1connect.NewUserServiceHandler(userHandler, interceptors))

	return &http.Server{Handler: mux}
}
