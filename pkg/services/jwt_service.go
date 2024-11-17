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
	ErrCreatingToken         = errors.New("error creating jwt token")
	ErrInvalidKidTokenHeader = errors.New("invalid or missing 'kid' in token header")
	ErrInvalidToken          = errors.New("invalid token")
)

var privKey *rsa.PrivateKey
var PubKey *rsa.PublicKey

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
	PubKey = loadPubKey

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

func GetClaims(token string, claims jwt.Claims) error {
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

		return PubKey, nil
	})

	if err != nil {
		return err
	}

	if !jwtToken.Valid {
		return ErrInvalidToken
	}

	return nil
}

func ValidateToken(token string) bool {
	return GetClaims(token, &AuthClaims{}) == nil
}
