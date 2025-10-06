package utils

import (
	"time"

	"github.com/BeatEcoprove/identityService/pkg/adapters"
	"github.com/stretchr/testify/mock"
)

type (
	MockRedis struct {
		mock.Mock
	}

	MockRabbitMq struct {
		mock.Mock
	}
)

func (r *MockRedis) GetValue(key adapters.RedisKey) (string, error) {
	args := r.Called(key)
	return args.String(0), args.Error(1)
}

func (r *MockRedis) SetValue(key adapters.RedisKey, value any, expiration time.Duration) error {
	args := r.Called(key, value, expiration)
	return args.Error(0)
}

func (r *MockRedis) GetAndDelValue(key adapters.RedisKey) (string, error) {
	args := r.Called(key)
	return args.String(0), args.Error(1)
}

func (r *MockRedis) Close() error {
	args := r.Called()
	return args.Error(0)
}

func (rc *MockRabbitMq) Publish(payload adapters.BrokerPayload, topic adapters.BrokerScope) error {
	args := rc.Called(payload)
	return args.Error(0)
}

func (rc *MockRabbitMq) Close() error {
	args := rc.Called()
	return args.Error(0)
}
