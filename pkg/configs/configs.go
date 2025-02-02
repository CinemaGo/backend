package configs

import (
	"fmt"
	"os"
)

// LoadEnvironmentVariable retrieves the value of an environment variable
// by its key. If the environment variable is not set or is empty, it returns
// an error indicating that the variable is missing.
//
// Parameters:
//
//	key (string): The name of the environment variable to retrieve.
//
// Returns:
//
//	string: The value of the environment variable if it exists.
//	error: An error if the environment variable is missing or empty.
func LoadEnvironmentVariable(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("config: missing environment variable: %s", key)
	}

	return value, nil
}
