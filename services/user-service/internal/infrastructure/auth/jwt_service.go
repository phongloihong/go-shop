package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	domain_error "github.com/phongloihong/go-shop/services/user-service/internal/domain/domain_errors"
	"github.com/phongloihong/go-shop/services/user-service/internal/domain/entity"
	"github.com/phongloihong/go-shop/services/user-service/internal/domain/service"
	"github.com/phongloihong/go-shop/services/user-service/internal/pkg/utils"
)

type JWTService struct {
	accessSecret     []byte
	refreshSecret    []byte
	accessExpiresIn  time.Duration
	refreshExpiresIn time.Duration
}

func NewJWTService(accessSecret []byte, refreshSecret []byte, accessExpiresIn, refreshExpiresIn time.Duration) service.AuthService {
	return &JWTService{
		accessSecret:     accessSecret,
		refreshSecret:    refreshSecret,
		accessExpiresIn:  accessExpiresIn,
		refreshExpiresIn: refreshExpiresIn,
	}
}

type customClaims struct {
	service.TokenClaims
	jwt.RegisteredClaims
}

func (j *JWTService) GenerateToken(user *entity.User) (*service.TokenPairs, error) {
	createTime := time.Now()

	accessToken, err := j.signToken(user, createTime, j.accessSecret, j.accessExpiresIn)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.signToken(user, createTime, j.refreshSecret, j.refreshExpiresIn)
	if err != nil {
		return nil, err
	}

	return &service.TokenPairs{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(j.accessExpiresIn.Seconds()),
	}, err
}

func (j *JWTService) ValidateToken(tokenString string, secret []byte) (*service.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (any, error) {
		// check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain_error.NewInvalidData(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}

		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		return &claims.TokenClaims, nil
	}

	return nil, domain_error.NewInvalidData("invalid token claims or token is not valid")
}

func (j *JWTService) signToken(user *entity.User, createTime time.Time, secret []byte, expiresIn time.Duration) (string, error) {
	claims := &customClaims{
		TokenClaims: service.TokenClaims{
			UserID: user.ID,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "UserService",
			Subject:   "",
			Audience:  jwt.ClaimStrings{},
			ExpiresAt: jwt.NewNumericDate(createTime.Add(expiresIn)),
			NotBefore: jwt.NewNumericDate(createTime),
			IssuedAt:  jwt.NewNumericDate(createTime),
			ID:        utils.NewUUID(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
