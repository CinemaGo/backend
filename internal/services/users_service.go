package services

import (
	"cinemaGo/backend/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	InsertNew(name, surname, email, phoneNumber, password string) error
	UserAuthentication(email, password string) (int, string, error)
	FetchUserInformations(userID int) (models.UserInfo, error)
	UpdateUserInformations(userID int, name, surname, phoneNumber string) error
}

type UserService struct {
	db          models.DBContractUsers
	redisClient *redis.Client
}

func NewUsersService(db models.DBContractUsers) (*UserService, error) {
	redisAddr, redisPass, err := LoadRedisEnvironmentVariables("REDIS_ADDR", "REDIS_PASS")
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:        redisAddr,
		Password:    redisPass,
		DB:          2,
		DialTimeout: 5 * time.Second,
	})
	return &UserService{
		db:          db,
		redisClient: rdb,
	}, nil
}

// InsertNew creates a new user by hashing the provided password and inserting the user data into the database.
//
// Parameters:
// - name: The first name of the user.
// - surname: The last name of the user.
// - email: The email address of the user.
// - phoneNumber: The phone number of the user.
// - password: The password provided by the user.
//
// Returns:
// - error: Returns an error if there is any issue with password hashing or inserting the user into the database.
// If the email is already taken, it returns ErrDuplicatedEmail.
func (us *UserService) InsertNew(name, surname, email, phoneNumber, password string) error {
	// Hash the provided password using bcrypt.
	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error occurred while hashing password in the service section: %w", err)
	}

	// Insert the new user into the database.
	err = us.db.InsertNewUser(name, surname, email, phoneNumber, password_hash)
	if err != nil {
		// Handle duplicate email error
		if errors.Is(err, models.ErrDuplicatedEmail) {
			return ErrDuplicatedEmail
		}
		return fmt.Errorf("error occurred while inserting new user in the service section: %w", err)
	}

	return nil
}

// UserAuthentication verifies a user's credentials by comparing the provided password with the stored hashed password.
//
// Parameters:
// - email: The email of the user attempting to authenticate.
// - password: The password provided by the user for authentication.
//
// Returns:
// - userID: The ID of the authenticated user.
// - userRole: The role of the authenticated user (e.g., admin, regular user).
// - error: Returns an error if authentication fails, such as wrong password or user not found.
func (s *UserService) UserAuthentication(email, password string) (int, string, error) {
	// Retrieve the user's credentials (ID, hashed password, and role) from the database.
	userID, password_hash, userRole, err := s.db.RetrieveUserCredentials(email)
	if err != nil {
		// If user is not found, return ErrUserNotFound
		if errors.Is(err, models.ErrUserNotFound) {
			return 0, "", ErrUserNotFound
		}
		return 0, "", fmt.Errorf("error occurred while user authentication in the service section: %v", err)
	}

	// Compare the provided password with the stored hashed password.
	if err := bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password)); err != nil {
		// Return error if passwords don't match
		return 0, "", ErrUserInvalidCredentials
	}

	return userID, userRole, nil
}

// FetchUserInformations retrieves a user's information from Redis or, if not found, from the database.
//
// Parameters:
// - userID: The ID of the user whose information is being fetched.
//
// Returns:
// - userInfo: A struct containing the user's name, surname, email, and phone number.
// - error: Returns an error if the user is not found in the database or an issue occurs while fetching the data.
func (us *UserService) FetchUserInformations(userID int) (models.UserInfo, error) {
	// Attempt to fetch the user information from Redis cache.
	userInfo, err := us.fetchUserInformationsFromRedisCache("userID")
	if errors.Is(err, ErrNoCachedDataFound) {
		// If not found in cache, retrieve from database.
		userInfo, err := us.db.RetrieveUserInfo(userID)
		if err != nil {
			// Handle the case where the user does not exist.
			if errors.Is(err, models.ErrUserNotFound) {
				return models.UserInfo{}, ErrUserNotFound
			}
			return models.UserInfo{}, fmt.Errorf("error occurred while fetching user information in the service section: %w", err)
		}

		// Cache the user information in Redis for future requests.
		if err := us.cacheUserInformationsInRedis(userID, userInfo); err != nil {
			return models.UserInfo{}, err
		}

		return userInfo, nil
	}
	if err != nil {
		// Return any other error that might occur during fetching from cache.
		return models.UserInfo{}, nil
	}

	// Return the user info from cache if successfully retrieved.
	return userInfo, nil
}

// UpdateUserInformations updates the user's information in the database and cache.
//
// Parameters:
// - userID: The ID of the user whose information is being updated.
// - name: The new first name of the user.
// - surname: The new last name of the user.
// - phoneNumber: The new phone number of the user.
//
// Returns:
// - error: Returns an error if there is an issue updating the user information in the database or cache.
func (us *UserService) UpdateUserInformations(userID int, name, surname, phoneNumber string) error {
	// Update the user's information in the database.
	err := us.db.UpdateUserInformationByID(userID, name, surname, phoneNumber)
	if err != nil {
		return fmt.Errorf("error occurred while updating the user information in the service section: %w", err)
	}

	// Retrieve the updated user information from the database.
	userInfo, err := us.db.RetrieveUserInfo(userID)
	if err != nil {
		// Handle case where user is not found.
		if errors.Is(err, models.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("error occurred while fetching user information in the service section: %w", err)
	}

	// Cache the updated user information in Redis.
	if err := us.cacheUserInformationsInRedis(userID, userInfo); err != nil {
		return err
	}

	return nil
}

// cacheUserInformationsInRedis stores the user's information in Redis cache.
//
// Parameters:
// - keyName: The Redis cache key for storing the user information.
// - userInfo: The user information to store in the cache.
//
// Returns:
// - error: Returns an error if there is an issue marshalling or storing the data in Redis.
func (us *UserService) cacheUserInformationsInRedis(keyName int, userInfo models.UserInfo) error {
	// Marshal the user information to JSON format.
	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		return fmt.Errorf("error occurred while marshalling userInfo, to cache for Redis: %w", err)
	}

	// Store the marshalled user information in Redis.
	err = us.redisClient.Set(context.Background(), fmt.Sprintf("%v", keyName), jsonData, 30*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("error occurred while setting up userInfo data in Redis: %w", err)
	}

	return nil
}

// fetchUserInformationsFromRedisCache attempts to retrieve the user's information from Redis cache.
//
// Parameters:
// - keyName: The Redis cache key used to retrieve the user information.
//
// Returns:
// - userInfo: The user information retrieved from Redis cache.
// - error: Returns ErrNoCachedDataFound if the user information is not found in cache.
func (us *UserService) fetchUserInformationsFromRedisCache(keyName string) (models.UserInfo, error) {
	// Get the user information from Redis.
	data, err := us.redisClient.Get(context.Background(), keyName).Result()
	if errors.Is(err, redis.Nil) {
		// Return an error if no cached data is found.
		return models.UserInfo{}, ErrNoCachedDataFound
	}

	if err != nil {
		return models.UserInfo{}, fmt.Errorf("error occurred while fetching cached userInfo data from Redis: %w", err)
	}

	var userInfo models.UserInfo

	// Unmarshal the cached data into the UserInfo struct.
	err = json.Unmarshal([]byte(data), &userInfo)
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("error occurred while unmarshalling userInfo data that is coming from Redis cache: %w", err)
	}

	// Return the user information if successfully retrieved from cache.
	return userInfo, nil
}
