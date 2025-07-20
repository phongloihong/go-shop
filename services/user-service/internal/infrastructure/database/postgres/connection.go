package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/phongloihong/go-shop/services/user-service/internal/config"
)

func NewConnection(ctx context.Context, cfg *config.DatabaseConfig) (*pgx.Conn, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
