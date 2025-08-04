package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

func isDuplicateKeyError(err error) bool {
	if pgxErr, ok := err.(*pgconn.PgError); ok {
		fmt.Println("pgxErr code:", pgxErr.Code)
		return pgxErr.Code == "23505" // Unique violation
	}

	return false
}
