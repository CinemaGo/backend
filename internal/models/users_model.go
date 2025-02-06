package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type DBContractUsers interface {
	InsertNewUser(name, surname, email, phoneNumber string, password_hash []byte) error
	RetrieveUserCredentials(email string) (int, string, string, error)
	RetrieveUserInfo(userID int) (UserInfo, error)
	UpdateUserInformationByID(userID int, name, surname, phoneNumber string) error
}

type Users struct {
	db DBContractUsers
}

func NewUsers(db DBContractUsers) *Users {
	return &Users{db: db}
}

// InsertNewUser inserts a new user into the database.
// If the email already exists, it returns ErrDuplicatedEmail.
//
// Parameters:
// - name: The user's first name.
// - surname: The user's last name.
// - email: The user's email address (must be unique).
// - phoneNumber: The user's phone number.
// - password_hash: The hashed password of the user.
//
// Returns:
// - nil if the user is inserted successfully.
// - ErrDuplicatedEmail if the email already exists in the database.
// - error if a database issue occurs.
func (psql *Postgres) InsertNewUser(name, surname, email, phoneNumber string, password_hash []byte) error {
	stmt := `INSERT INTO users (name, surname, email, phone_number, password_hash) VALUES ($1, $2, $3, $4, $5)`

	// Execute the query, passing the parameters to the query placeholders.
	_, err := psql.DB.Exec(stmt, name, surname, email, phoneNumber, password_hash)
	if err != nil {
		// Check if the error is related to a duplicate email (PostgreSQL constraint violation).
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			// Return a custom error for duplicated email.
			return ErrDuplicatedEmail
		}
		// Return a generic error with the original database error wrapped.
		return fmt.Errorf("failed to insert new user into the database: %w", err)
	}
	return nil
}

// RetrieveUserCredentials retrieves the user's credentials (ID, password hash, and role) based on their email.
// Returns ErrUserNotFound if the user with the specified email does not exist.
//
// Parameters:
// - email: The email address of the user whose credentials are to be fetched.
//
// Returns:
// - userID: The unique identifier of the user.
// - password_hash: The user's hashed password.
// - userRole: The user's role (e.g., admin, user).
// - error if a database issue occurs or the user is not found.
func (psql *Postgres) RetrieveUserCredentials(email string) (int, string, string, error) {
	stmt := `SELECT id, password_hash, role FROM users WHERE email = $1`
	var userID int
	var password_hash string
	var userRole string

	// Execute the query and scan the results into the respective variables.
	err := psql.DB.QueryRow(stmt, email).Scan(&userID, &password_hash, &userRole)
	if err != nil {
		// If no user is found, return the custom error ErrUserNotFound.
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", "", ErrUserNotFound
		}
		// Return a generic error with the original database error wrapped.
		return 0, "", "", fmt.Errorf("failed to retrieve user credentials: %w", err)
	}

	return userID, password_hash, userRole, nil
}

// RetrieveUserInfo retrieves detailed information (name, surname, email, and phone number) of a user based on their userID.
// Returns ErrUserNotFound if the user with the specified ID does not exist.
//
// Parameters:
// - userID: The unique identifier of the user.
//
// Returns:
// - user: A UserInfo struct containing the user's name, surname, email, and phone number.
// - error if a database issue occurs or the user is not found.
func (psql *Postgres) RetrieveUserInfo(userID int) (UserInfo, error) {
	stmt := `SELECT name, surname, email, phone_number FROM users WHERE id = $1`

	var user UserInfo

	// Execute the query and scan the results into the 'user' struct.
	err := psql.DB.QueryRow(stmt, userID).Scan(&user.Name, &user.Surname, &user.Email, &user.PhoneNumber)
	if err != nil {

		// If no user is found, return the custom error ErrUserNotFound.
		if errors.Is(err, sql.ErrNoRows) {
			return UserInfo{}, ErrUserNotFound
		}

		// Return a generic error with the original database error wrapped.
		return UserInfo{}, fmt.Errorf("failed to retrieve user information: %w", err)
	}

	return user, nil
}

// UpdateUserInformationByID updates the user's information (name, surname, and phone number) based on their userID.
// The updated timestamp is automatically set to the current time.
//
// Parameters:
// - userID: The unique identifier of the user.
// - name: The new first name of the user.
// - surname: The new last name of the user.
// - phoneNumber: The new phone number of the user.
//
// Returns:
// - nil if the user information is updated successfully.
// - error if a database issue occurs during the update.
func (psql *Postgres) UpdateUserInformationByID(userID int, name, surname, phoneNumber string) error {
	stmt := `UPDATE users SET name = $1, surname = $2, phone_number = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`

	// Execute the query to update the user information.
	_, err := psql.DB.Exec(stmt, name, surname, phoneNumber, userID)
	if err != nil {
		// Return a generic error with the original database error wrapped.
		return fmt.Errorf("failed to update user information: %w", err)
	}

	return nil
}
