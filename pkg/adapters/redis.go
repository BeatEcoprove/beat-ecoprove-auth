package adapters

import (
	"strings"
	"time"
)

const Delimiter = ":"

type (
	Redis interface {
		GetValue(key RedisKey) (string, error)
		SetValue(key RedisKey, value interface{}, expiration time.Duration) error
		GetAndDelValue(key RedisKey) (string, error)
		Close() error
	}

	RedisKey struct {
		Values []string
		Key    string
	}
)

func NewRedisKey(values ...string) RedisKey {
	var key strings.Builder

	for i := 0; i < len(values); i++ {
		key.WriteString(values[i])

		if i < len(values)-1 {
			key.WriteString(Delimiter)
		}
	}

	return RedisKey{
		Key:    key.String(),
		Values: values,
	}
}
