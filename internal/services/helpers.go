package services

import (
	"cinemaGo/backend/pkg/configs"
	"fmt"
)

func LoadRedisEnvironmentVariables(address, password string) (string, string, error) {
	redisAddr, err := configs.LoadEnvironmentVariable(address)
	if err != nil {
		return "", "", fmt.Errorf("failed to receive redis address from .env file")
	}

	redisPass, err := configs.LoadEnvironmentVariable(password)
	if err != nil {
		return "", "", fmt.Errorf("failed to receive redis password from .env file")
	}

	return redisAddr, redisPass, nil
}
