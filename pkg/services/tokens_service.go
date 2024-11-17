package services

import (
	"time"

	"github.com/BeatEcoprove/identityService/config"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
)

type (
	ITokenService interface {
		CreateAuthenticationTokens(payload TokenPayload) (*JwtToken, *JwtToken, error)
		ValidateToken(authId, token string) error
	}

	TokenService struct {
		redis interfaces.Redis
	}
)

func NewTokenService(redis interfaces.Redis) *TokenService {
	return &TokenService{
		redis: redis,
	}
}

func NewAccessTokenKey(userId string) interfaces.RedisKey {
	return interfaces.NewRedisKey(userId, "access")
}

func NewRefreshTokenKey(userId string) interfaces.RedisKey {
	return interfaces.NewRedisKey(userId, "refresh")
}

func generateAuthenticationTokens(payload TokenPayload, accessTokenExp, refreshTokenExp time.Duration) (*JwtToken, *JwtToken, error) {
	payload.Duration = time.Now().Add(accessTokenExp)
	accessToken, err := CreateJwtToken(payload)

	if err != nil {
		return nil, nil, err
	}

	payload.Duration = time.Now().Add(refreshTokenExp)
	refreshToken, err := CreateJwtToken(payload)

	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (ts *TokenService) ValidateToken(authId, token string) error {
	storedToken, err := ts.redis.GetValue(NewAccessTokenKey(authId))

	if err != nil {
		return err
	}

	if storedToken != token {
		return ErrInvalidToken
	}

	return nil
}

func (ts *TokenService) CreateAuthenticationTokens(payload TokenPayload) (*JwtToken, *JwtToken, error) {
	env := config.GetCofig()

	accessTokenExp := time.Duration(env.JWT_ACCESS_EXPIRED) * time.Minute           // per minute
	refreshTokenExp := time.Duration(env.JWT_REFRESH_EXPIRED) * time.Hour * 24 * 30 // per month

	accessToken, refreshToken, err := generateAuthenticationTokens(payload, accessTokenExp, refreshTokenExp)

	if err != nil {
		return nil, nil, err
	}

	ts.redis.GetAndDelValue(NewAccessTokenKey(payload.UserId))
	ts.redis.GetAndDelValue(NewRefreshTokenKey(payload.UserId))

	if err := ts.redis.SetValue(NewAccessTokenKey(payload.UserId), accessToken.Token, accessTokenExp); err != nil {
		return nil, nil, ErrCreatingToken
	}

	if err := ts.redis.SetValue(NewRefreshTokenKey(payload.UserId), refreshToken.Token, refreshTokenExp); err != nil {
		return nil, nil, ErrCreatingToken
	}

	return accessToken, refreshToken, nil
}
