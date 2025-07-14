package postgres

import (
	"context"
	"travel-planning/internal/domain/entity"
	"travel-planning/internal/infrastructure/database/postgres/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(db sqlc.DBTX) *UserRepository {
	return &UserRepository{
		queries: sqlc.New(db),
	}
}

func (ur *UserRepository) CreateUser(user *entity.User) (*entity.User, error) {
	return nil, nil
}

func (*UserRepository) UpdateUser(user *entity.User) (*entity.User, error) {
	return nil, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	uuid := pgtype.UUID{}
	uuid.Scan(id)

	user, err := ur.queries.GetUserByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return entity.UserFromDatabase(
		user.ID.String(),
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone.String,
		user.Password,
		user.CreatedAt.Time.Unix(),
		user.UpdatedAt.Time.Unix(),
	), nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := ur.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return entity.UserFromDatabase(
		user.ID.String(),
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone.String,
		user.Password,
		user.CreatedAt.Time.Unix(),
		user.UpdatedAt.Time.Unix(),
	), err
}

func (ur *UserRepository) GetPublicProfileByIds(ctx context.Context, ids []string) ([]*entity.UserPublicProfile, error) {
	ret := make([]*entity.UserPublicProfile, 0)
	users, err := ur.queries.GetPublicProfileByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		ret = append(ret, entity.NewUserPublicProfile(
			user.ID.String(),
			user.FirstName,
			user.LastName,
		))
	}

	return ret, err
}
