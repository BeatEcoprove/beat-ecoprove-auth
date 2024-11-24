package utils

import (
	"encoding/pem"
	"errors"
	"path"

	"github.com/BeatEcoprove/identityService/config"
	"github.com/BeatEcoprove/identityService/pkg/services"
)

var (
	ErrKillContainer = errors.New("failed to terminate container")
)

func GetDotEnvPath() (string, error) {
	projectPath, err := config.GetProjectRoot()

	if err != nil {
		return "", err
	}

	return path.Join(projectPath, config.DotEnv), nil
}

func TestSetup() {
	// before
	dotEnvPath, err := GetDotEnvPath()

	if err != nil {
		panic(err)
	}

	config.LoadEnv(dotEnvPath)

	// load pki
	keys, err := services.CreatePKI()

	if err != nil {
		panic(err)
	}

	privateKey := pem.EncodeToMemory(keys.PrivateKey)
	publicKey := pem.EncodeToMemory(keys.PublicKey)

	err = services.LoadKeys(publicKey, privateKey)

	if err != nil {
		panic(err)
	}

	if _, err = services.CreateJWKS(); err != nil {
		panic(err)
	}
}
