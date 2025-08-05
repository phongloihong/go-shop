package service

import (
	"github.com/phongloihong/go-shop/services/user-service/internal/domain/entity"
)

type (
	TokenClaims struct {
		UserID  string
		TokenID string
	}

	TokenPairs struct {
		AccessToken  string
		RefreshToken string
		ExpiresIn    int64
	}
)

type AuthService interface {
	// Token life cycle management
	GenerateToken(user *entity.User) (*TokenPairs, error)
	ValidateToken(token string, secret []byte) (*TokenClaims, error)
	// RefreshToken(token string) (*TokenPairs, error)

	// Token Management
	// RevokeToken(tokenID string) error
	// IsTokenRevoked(tokenID string) bool
}
