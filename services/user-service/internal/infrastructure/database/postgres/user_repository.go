package postgres

import (
	"context"
	"fmt"
	"time"

	domain_error "github.com/phongloihong/go-shop/services/user-service/internal/domain/domain_errors"
	"github.com/phongloihong/go-shop/services/user-service/internal/domain/entity"
	"github.com/phongloihong/go-shop/services/user-service/internal/infrastructure/database/postgres/sqlc"

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
	if err := phone.Scan(user.Phone.String()); err != nil {
		return nil, domain_error.NewInvalidData(fmt.Sprintf("invalid phone number: %s", user.Phone.String()))
	}

	timeNow := pgtype.Timestamp{}
	if err := timeNow.Scan(time.Now()); err != nil {
		return nil, domain_error.NewInvalidData(fmt.Sprintf("failed to scan current time: %s", err.Error()))
	}

	hashPassword, err := user.Password.Hash()
	if err != nil {
		return nil, domain_error.NewInternalError(fmt.Sprintf("failed to hash password: %s", err.Error()))
	}

	newUser, err := ur.queries.InsertUser(ctx, sqlc.InsertUserParams{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email.String(),
		Phone:     phone,
		Password:  hashPassword,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, domain_error.NewAlreadyExistsError(fmt.Sprintf("user with email %s already exists", user.Email.String()))
		}

		return nil, domain_error.NewInternalError(fmt.Sprintf("failed to create user: %s", err.Error()))
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
		return 0, domain_error.NewInvalidData(fmt.Sprintf("invalid user ID: %s", user.ID))
	}

	phone := pgtype.Text{}
	if err := phone.Scan(user.Phone); err != nil {
		return 0, domain_error.NewInvalidData(fmt.Sprintf("invalid phone number: %s", user.Phone))
	}

	updatedAt := pgtype.Timestamp{}
	if err := updatedAt.Scan(user.UpdatedAt.Time()); err != nil {
		return 0, domain_error.NewInvalidData(fmt.Sprintf("failed to scan updated timestamp: %s", err.Error()))
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
		return 0, domain_error.NewInternalError(fmt.Sprintf("failed to update user: %s", err.Error()))
	}

	return ret.RowsAffected(), nil
}

func (ur *UserRepository) ChangePassword(ctx context.Context, id string, newPassword string) (int64, error) {
	uuid := pgtype.UUID{}
	if err := uuid.Scan(id); err != nil {
		return 0, domain_error.NewInvalidData(fmt.Sprintf("invalid user ID: %s", id))
	}

	updatedAt := pgtype.Timestamp{}
	if err := updatedAt.Scan(time.Now()); err != nil {
		return 0, domain_error.NewInvalidData(fmt.Sprintf("failed to scan current time: %s", err.Error()))
	}

	updateParams := sqlc.UpdateUserPasswordParams{
		ID:        uuid,
		Password:  newPassword,
		UpdatedAt: updatedAt,
	}
	ret, err := ur.queries.UpdateUserPassword(ctx, updateParams)
	if err != nil {
		return 0, domain_error.NewInternalError(fmt.Sprintf("failed to change password: %s", err.Error()))
	}

	return ret.RowsAffected(), nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	uuid := pgtype.UUID{}
	if err := uuid.Scan(id); err != nil {
		return nil, domain_error.NewInvalidData(fmt.Sprintf("invalid user ID: %s", id))
	}

	user, err := ur.queries.GetUserByID(ctx, uuid)
	if err != nil {
		return nil, domain_error.NewInternalError(fmt.Sprintf("failed to get user by ID: %s", err.Error()))
	}

	return ur.sqlcUserToEntity(user), nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := ur.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, domain_error.NewInternalError(fmt.Sprintf("failed to get user by email: %s", err.Error()))
	}

	return ur.sqlcUserToEntity(user), nil
}

func (ur *UserRepository) GetPublicProfileByIds(ctx context.Context, ids []string) ([]*entity.UserPublicProfile, error) {
	ret := make([]*entity.UserPublicProfile, 0)
	users, err := ur.queries.GetPublicProfileByIds(ctx, ids)
	if err != nil {
		return nil, domain_error.NewInternalError(fmt.Sprintf("failed to get public profiles by IDs: %s", err.Error()))
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
