package services

import (
	"log"
	"time"

	"github.com/BeatEcoprove/identityService/config"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
)

type (
	TokenKey string

	ITokenService interface {
		CreateAuthenticationTokens(payload TokenPayload) (*JwtToken, *JwtToken, error)
		ValidateToken(authID, token string, key TokenKey) error
	}

	TokenService struct {
		redis interfaces.Redis
	}
)

const (
	AccessTokenKey  TokenKey = "access"
	RefreshTokenKey TokenKey = "refresh"
)

func NewTokenService(redis interfaces.Redis) *TokenService {
	return &TokenService{
		redis: redis,
	}
}

func NewAccessTokenKey(userId string) interfaces.RedisKey {
	return interfaces.NewRedisKey(userId, string(AccessTokenKey))
}

func NewRefreshTokenKey(userId string) interfaces.RedisKey {
	return interfaces.NewRedisKey(userId, string(RefreshTokenKey))
}

func generateAuthenticationTokens(payload TokenPayload, accessTokenExp, refreshTokenExp time.Duration) (*JwtToken, *JwtToken, error) {
	payload.Duration = time.Now().Add(accessTokenExp)
	payload.Type = Access

	accessToken, err := CreateJwtToken(payload)

	if err != nil {
		return nil, nil, err
	}

	payload.Duration = time.Now().Add(refreshTokenExp)
	payload.Type = Refresh

	refreshToken, err := CreateJwtToken(payload)

	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (ts *TokenService) ValidateToken(authId, token string, key TokenKey) error {
	var tokenKey interfaces.RedisKey

	switch key {
	case AccessTokenKey:
		tokenKey = NewAccessTokenKey(authId)
	case RefreshTokenKey:
		tokenKey = NewRefreshTokenKey(authId)
	default:
		return ErrInvalidToken
	}

	storedToken, err := ts.redis.GetValue(tokenKey)

	if err != nil {
		return err
	}

	if storedToken != token {
		return ErrInvalidToken
	}

	return nil
}

func (ts *TokenService) CreateAuthenticationTokens(payload TokenPayload) (*JwtToken, *JwtToken, error) {
	env := config.GetConfig()

	accessTokenExp := time.Duration(env.JWT_ACCESS_EXPIRED) * time.Minute           // per minute
	refreshTokenExp := time.Duration(env.JWT_REFRESH_EXPIRED) * time.Hour * 24 * 30 // per month

	accessToken, refreshToken, err := generateAuthenticationTokens(payload, accessTokenExp, refreshTokenExp)

	if err != nil {
		return nil, nil, err
	}

	ts.redis.GetAndDelValue(NewAccessTokenKey(payload.UserID))
	ts.redis.GetAndDelValue(NewRefreshTokenKey(payload.UserID))

	if err := ts.redis.SetValue(NewAccessTokenKey(payload.UserID), accessToken.Token, accessTokenExp); err != nil {
		log.Printf("%s", err.Error())
		return nil, nil, ErrCreatingToken
	}

	if err := ts.redis.SetValue(NewRefreshTokenKey(payload.UserID), refreshToken.Token, refreshTokenExp); err != nil {
		log.Printf("%s", err.Error())
		return nil, nil, ErrCreatingToken
	}

	return accessToken, refreshToken, nil
}
