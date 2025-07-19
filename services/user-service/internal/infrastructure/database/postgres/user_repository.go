package postgres

import (
	"context"
	"fmt"
	"time"
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

func (ur *UserRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	phone := pgtype.Text{}
	if err := phone.Scan(user.Phone); err != nil {
		return nil, fmt.Errorf("failed to scan phone: %w", err)
	}

	timeNow := pgtype.Timestamp{}
	if err := timeNow.Scan(time.Now()); err != nil {
		return nil, fmt.Errorf("failed to scan timestamp: %w", err)
	}

	newUser, err := ur.queries.InsertUser(ctx, sqlc.InsertUserParams{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email.String(),
		Phone:     phone,
		Password:  user.Password.String(),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	ret := entity.UserFromDatabase(
		newUser.ID.String(),
		newUser.FirstName,
		newUser.LastName,
		newUser.Email,
		newUser.Phone.String,
		newUser.Password,
		newUser.CreatedAt.Time.Unix(),
		newUser.UpdatedAt.Time.Unix(),
	)

	return ret, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, user *entity.User) (int64, error) {
	uuid := pgtype.UUID{}
	if err := uuid.Scan(user.ID); err != nil {
		return 0, fmt.Errorf("failed to scan user ID: %w", err)
	}

	phone := pgtype.Text{}
	if err := phone.Scan(user.Phone); err != nil {
		return 0, fmt.Errorf("failed to scan user phone: %w", err)
	}

	updatedAt := pgtype.Timestamp{}
	if err := updatedAt.Scan(user.UpdatedAt.Time()); err != nil {
		return 0, fmt.Errorf("failed to scan updated timestamp: %w", err)
	}

	updateParams := sqlc.UpdateUserParams{
		ID:        uuid,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email.String(),
		Phone:     phone,
		UpdatedAt: updatedAt,
	}
	ret, err := ur.queries.UpdateUser(ctx, updateParams)
	if err != nil {
		return 0, fmt.Errorf("failed to update user: %w", err)
	}

	return ret.RowsAffected(), nil
}

func (ur *UserRepository) ChangePassword(ctx context.Context, id string, newPassword string) (int64, error) {
	uuid := pgtype.UUID{}
	if err := uuid.Scan(id); err != nil {
		return 0, fmt.Errorf("failed to scan user ID: %w", err)
	}

	updatedAt := pgtype.Timestamp{}
	if err := updatedAt.Scan(time.Now()); err != nil {
		return 0, fmt.Errorf("failed to scan current time: %w", err)
	}

	updateParams := sqlc.UpdateUserPasswordParams{
		ID:        uuid,
		Password:  newPassword,
		UpdatedAt: updatedAt,
	}
	ret, err := ur.queries.UpdateUserPassword(ctx, updateParams)
	if err != nil {
		return 0, fmt.Errorf("failed to change password: %w", err)
	}

	return ret.RowsAffected(), nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	uuid := pgtype.UUID{}
	if err := uuid.Scan(id); err != nil {
		return nil, fmt.Errorf("failed to scan UUID: %w", err)
	}

	user, err := ur.queries.GetUserByID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return ur.sqlcUserToEntity(user), nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := ur.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return ur.sqlcUserToEntity(user), nil
}

func (ur *UserRepository) GetPublicProfileByIds(ctx context.Context, ids []string) ([]*entity.UserPublicProfile, error) {
	ret := make([]*entity.UserPublicProfile, 0)
	users, err := ur.queries.GetPublicProfileByIds(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get public profiles by IDs: %w", err)
	}

	for _, user := range users {
		ret = append(ret, entity.NewUserPublicProfile(
			user.ID.String(),
			user.FirstName,
			user.LastName,
		))
	}

	return ret, nil
}

func (*UserRepository) sqlcUserToEntity(sqlcUser sqlc.User) *entity.User {
	return entity.UserFromDatabase(
		sqlcUser.ID.String(),
		sqlcUser.FirstName,
		sqlcUser.LastName,
		sqlcUser.Email,
		sqlcUser.Phone.String,
		sqlcUser.Password,
		sqlcUser.CreatedAt.Time.Unix(),
		sqlcUser.UpdatedAt.Time.Unix(),
	)
}
