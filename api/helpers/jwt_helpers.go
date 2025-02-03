package helpers

import (
	"cinemaGo/backend/pkg/configs"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CreateAndSignJWT generates and signs a new JWT token for a user.
//
// Parameters:
// - userID: The ID of the user for whom the JWT is being created.
// - userRole: The role of the user (e.g., "admin", "user") to be included in the JWT claims.
//
// Returns:
// - A string representing the signed JWT.
// - An error if there's an issue creating or signing the token (e.g., missing secret key).
func CreateAndSignJWT(userID int, userRole string) (string, error) {
	// Create a new JWT with claims (userID, userRole, and ttl).
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   userID,                                // User ID included in the payload
		"userRole": userRole,                              // User Role included in the payload
		"ttl":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (24 hours)
	})

	// Load the JWT secret key from environment variables.
	secretKey, err := configs.LoadEnvironmentVariable("JWT_SECRET_KEY")
	if err != nil {
		// Return an error if the secret key cannot be retrieved from the environment.
		return "", fmt.Errorf("cannot get secret key while creating and signing JWT: %v", err)
	}

	// Sign the token with the secret key and return the token string.
	return token.SignedString([]byte(secretKey))
}

// SetCookie sets a secure cookie with the JWT token for client authentication.
//
// Parameters:
// - c: The Gin context, used to set the cookie in the response.
// - token: The JWT token to be included in the cookie.
//
// This function sets a cookie named "u_auth" with the JWT, which will be used for subsequent requests.
// The cookie has a 1-week expiration time and the SameSite attribute is set to Lax to prevent CSRF attacks.
func SetCookie(c *gin.Context, token string) {
	// Set the SameSite attribute for the cookie to Lax, preventing CSRF attacks.
	c.SetSameSite(http.SameSiteLaxMode)

	// Set the cookie with the JWT token, with a duration of 1 week (604800 seconds).
	c.SetCookie("u_auth", token, 604800, "", "", false, true)
}

// UnauthorizedResponse sends a 401 Unauthorized response with a login prompt message.
// It is used when authentication fails, aborting the request to prevent further processing.
//
// Parameters:
// - c: The Gin context used to send the response.
//
// This function returns a JSON response with a message asking the user to log in to continue,
// and it halts further processing of the request using c.Abort().
func UnauthorizedResponse(c *gin.Context) {
	// Send a 401 Unauthorized response with a message prompting the user to log in.
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "Please log in to continue",
	})

	// Abort the request to prevent further processing.
	c.Abort()
}
