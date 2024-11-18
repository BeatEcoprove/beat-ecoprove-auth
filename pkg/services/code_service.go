package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
)

type (
	Code string

	CodeKey struct {
		Value string
		Exp   time.Duration
	}

	IPGService interface {
		CreateAndStoreCode(userId string) (*Code, error)
		ValidateCode(userId, code string) error
	}

	PGService struct {
		redis interfaces.Redis
	}
)

const (
	MinPassWordGen = 32
	Delimiter      = ":"
)

var (
	ForgotCodeKey = CodeKey{
		Value: "forgot",
		Exp:   15 * time.Minute, // 15 min
	}

	ErrCodeNotValid = errors.New("code is not valid")
)

func NewPGService(redis interfaces.Redis) *PGService {
	return &PGService{
		redis: redis,
	}
}

func NewForgotCodeKey(userId string) interfaces.RedisKey {
	return interfaces.NewRedisKey(userId, ForgotCodeKey.Value)
}

func (pgs *PGService) ValidateCode(userId, code string) error {
	var forgotKey = NewForgotCodeKey(userId)

	storedCode, err := pgs.redis.GetAndDelValue(forgotKey)

	if err != nil {
		return err
	}

	base64Content := strings.Split(storedCode, Delimiter)

	if len(base64Content) > 2 {
		return ErrCodeNotValid
	}

	codeKey, err := base64.StdEncoding.DecodeString(base64Content[1])

	if err != nil {
		return err
	}

	rawCode, err := AesDecrypt(base64Content[0], codeKey)

	if err != nil {
		return err
	}

	if rawCode != code {
		return ErrCodeNotValid
	}

	return nil
}

func (pgs *PGService) CreateAndStoreCode(userId string) (*Code, error) {
	var forgotKey = NewForgotCodeKey(userId)
	code, err := GenerateCode()

	if err != nil {
		return nil, err
	}

	pgs.redis.GetAndDelValue(forgotKey)

	genKey, err := GeneratePassword(MinPassWordGen, MinPassWordGen)

	if err != nil {
		return nil, err
	}

	cipherCode, err := AesEncrypt([]byte(code), []byte(genKey))

	if err != nil {
		return nil, err
	}

	payload := fmt.Sprintf("%s%s%s", cipherCode, Delimiter, base64.StdEncoding.EncodeToString([]byte(genKey)))
	if err := pgs.redis.SetValue(forgotKey, payload, ForgotCodeKey.Exp); err != nil {
		return nil, ErrCreatingToken
	}

	return (*Code)(&code), nil
}
