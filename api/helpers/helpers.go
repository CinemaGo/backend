package helpers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServerError handles unexpected server-side errors and returns a generic response to the client.
// It logs the error for server-side troubleshooting while ensuring that sensitive error details
// are not exposed to the client.
//
// Parameters:
//
//	c (*gin.Context): The Gin context used to build the response.
//	err (error): The error that occurred during server processing.
//
// Returns:
//
//		void: Responds directly to the client with a 500 Internal Server Error status and a
//	      general message indicating a server issue.
func ServerError(c *gin.Context, err error) {
	// Log the server-side error (don't expose to the user)
	log.Println(err)
	// Send a generic error response to the client
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "An unexpected server error occurred. Please try again later.",
	})
}

// ClientError handles client-side errors and returns an appropriate HTTP response.
// It allows for custom error messages and status codes to be sent to the client based on
// the type of error that occurred.
//
// Parameters:
//
//	c (*gin.Context): The Gin context used to build the response.
//	code (int): The HTTP status code that indicates the type of error (e.g., 400, 404, etc.)
//	message (string): A descriptive error message that explains the client-side issue.
//
// Returns:
//
//	void: Responds directly to the client with the provided error message and status code.
func ClientError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"error": message, // Key "error" is a good convention
	})
}
