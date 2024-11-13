package services

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/BeatEcoprove/identityService/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type (
	JwtToken struct {
		Token    string
		ExpireAt int
	}

	TokenPayload struct {
		Email      string
		UserId     string
		ProfileId  string
		ProfileIds []string
		Role       string
		Duration   time.Time
	}

	AuthClaims struct {
		jwt.RegisteredClaims
		Email      string   `json:"email,omitempty"`
		Role       string   `json:"role,omitempty"`
		ProfileId  string   `json:"profileId,omitempty"`
		ProfileIds []string `json:"profileIds"`
	}
)

var (
	ErrInvalidKidTokenHeader = errors.New("invalid or missing 'kid' in token header")
	ErrInvalidToken          = errors.New("invalid token")
)

var privKey *rsa.PrivateKey
var pubKey *rsa.PublicKey

func LoadKeys(publicKey, privateKey []byte) error {
	loadPrivKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)

	if err != nil {
		return err
	}

	loadPubKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)

	if err != nil {
		return err
	}

	privKey = loadPrivKey
	pubKey = loadPubKey

	return nil
}

func CreateJwtToken(payload TokenPayload) (*JwtToken, error) {
	env := config.GetCofig()

	claims := AuthClaims{
		Email:      payload.Email,
		Role:       payload.Role,
		ProfileId:  payload.ProfileId,
		ProfileIds: payload.ProfileIds,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.JWT_ISSUER,
			Audience:  jwt.ClaimStrings{env.JWT_AUDIENCE},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ExpiresAt: &jwt.NumericDate{Time: payload.Duration},
			Subject:   payload.UserId,
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = generateKid()

	jwtToken, err := token.SignedString(privKey)

	if err != nil {
		return nil, err
	}

	return &JwtToken{
		Token:    jwtToken,
		ExpireAt: int(payload.Duration.UTC().Unix()),
	}, nil
}

func CreateAuthenticationTokens(payload TokenPayload) (*JwtToken, *JwtToken, error) {
	env := config.GetCofig()

	payload.Duration = time.Now().Add(time.Duration(env.JWT_ACCESS_EXPIRED) * time.Minute) // per minute
	accessToken, err := CreateJwtToken(payload)

	if err != nil {
		return nil, nil, err
	}

	payload.Duration = time.Now().Add(time.Duration(env.JWT_REFRESH_EXPIRED) * time.Hour * 24 * 30) // per month
	refreshToken, err := CreateJwtToken(payload)

	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func GetClaims(token string, claims jwt.Claims) (interface{}, error) {
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		kid, ok := t.Header["kid"].(string)

		if !ok {
			return nil, ErrInvalidKidTokenHeader
		}

		if storedJwks.Keys[0].Kid != kid {
			return nil, ErrInvalidToken
		}

		return pubKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func ValidateToken(token string) bool {
	_, err := GetClaims(token, &AuthClaims{})

	return err == nil
}
