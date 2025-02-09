package middlewares

import (
	"cinemaGo/backend/api/helpers"
	"cinemaGo/backend/pkg/configs"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// UserAuthorizationJWT is a Gin middleware function that checks the validity of a JWT token
// passed via a cookie ("u_auth"). It verifies the token's integrity, expiration, and claims
// before allowing the request to proceed. If the token is invalid or missing, it returns an
// Unauthorized response.
//
// Returns:
// - A gin.HandlerFunc which handles JWT verification and authorization.
func UserAuthorizationJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the JWT token from the cookie "u_auth"
		tokenStr, err := c.Cookie("u_auth")
		if err != nil {
			// If the token is missing or invalid, send an Unauthorized response
			helpers.UnauthorizedResponse(c)
			return
		}

		// Parse and validate the JWT token using the HMAC signing method
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token uses the correct signing method (HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// If the signing method is not HMAC, return an error
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Load the secret key for signing the JWT from environment variables
			secretKey, err := configs.LoadEnvironmentVariable("JWT_SECRET_KEY")
			if err != nil {
				// If the secret key can't be loaded, return an error
				return nil, fmt.Errorf("cannot get secret key while creating and signing JWT: %v", err)
			}

			// Return the secret key to verify the JWT signature
			return []byte(secretKey), nil
		})

		// If token parsing or validation fails, send an Unauthorized response
		if err != nil || !token.Valid {
			helpers.UnauthorizedResponse(c)
			return
		}

		// Extract the claims from the JWT token and check if they are valid
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			// If claims extraction fails, return an Unauthorized response
			helpers.UnauthorizedResponse(c)
			return
		}

		// Check if the token has expired based on the "ttl" claim (time-to-live)
		if claims["ttl"].(float64) < float64(time.Now().Unix()) {
			// If the token is expired, return an Unauthorized response
			helpers.UnauthorizedResponse(c)
			return
		}

		// Ensure the "userID" claim is valid, otherwise return Unauthorized
		userID := claims["userID"].(float64)
		if userID == 0 {
			// If userID is invalid, return an Unauthorized response
			helpers.UnauthorizedResponse(c)
			return
		}

		// Extract the "userRole" claim from the JWT token
		userRole := claims["userRole"].(string)

		// Set the "userID" and "userRole" in the context for further use by downstream handlers
		c.Set("userID", userID)     // Store userID in context.
		c.Set("userRole", userRole) // Store userRole in context.

		// Proceed to the next middleware or handler in the chain
		c.Next()
	}
}

// AdminRoleRequired checks if the user is an admin. If not, it responds with a 403 status code.
func AdminRoleRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")

		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Access forbidden: Admins only",
			})
			c.Abort() // Stop further processing
			return
		}
		c.Next() // Continue to the next handler if the user is an admin
	}
}
