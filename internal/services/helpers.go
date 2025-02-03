package services

import (
	"cinemaGo/backend/pkg/configs"
	"fmt"
)

// LoadRedisEnvironmentVariables loads Redis environment variables from a .env file.
// It retrieves the Redis address and password from the environment and returns them.
//
// Parameters:
// - address: The key name of the Redis address environment variable in the .env file.
// - password: The key name of the Redis password environment variable in the .env file.
//
// Returns:
// - redisAddr: The Redis server address (from the .env file).
// - redisPass: The Redis server password (from the .env file).
// - error: Returns an error if either environment variable is not found or if the values cannot be loaded.
func LoadRedisEnvironmentVariables(address, password string) (string, string, error) {
	// Load the Redis address from the environment using the provided address key.
	redisAddr, err := configs.LoadEnvironmentVariable(address)
	if err != nil {
		// Return an error if the Redis address is not found or fails to load.
		return "", "", fmt.Errorf("failed to receive redis address from .env file")
	}

	// Load the Redis password from the environment using the provided password key.
	redisPass, err := configs.LoadEnvironmentVariable(password)
	if err != nil {
		// Return an error if the Redis password is not found or fails to load.
		return "", "", fmt.Errorf("failed to receive redis password from .env file")
	}

	// Return the loaded Redis address and password if both are successfully retrieved.
	return redisAddr, redisPass, nil
}
