package helpers

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

// RespondWithValidationErrors processes and returns validation errors when the request
// body fails to meet the defined validation criteria. It formats the errors in a way that
// is easy for the client to understand, displaying the invalid fields and the respective
// validation tags.
//
// Parameters:
//
//		err (error): The validation error that contains details of the failed validations.
//		structType (interface{}): The struct type used for validation that will be used to
//	                            map the error fields and generate proper messages.
//
// Returns:
//
//		void: Responds directly to the client with a 400 Bad Request status and a list of
//	      validation error messages.
func RespondWithValidationErrors(c *gin.Context, err error, structType interface{}) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid validation error format, please make sure all requested formats are correct",
		})
		return
	}

	var errorMessages []string

	// Reflect over the struct to get fields and generate appropriate error messages
	structValue := reflect.TypeOf(structType)

	// Iterate over validation errors and append formatted error messages
	for _, err := range errs {
		field, _ := structValue.FieldByName(err.Field())

		fieldName := err.Field()

		// If a JSON tag exists for the field, use it instead of the field name
		if tag := field.Tag.Get("json"); tag != "" {
			fieldName = tag
		}
		// Append detailed error message (field name and validation tag)
		errorMessages = append(errorMessages, fmt.Sprintf("'%s' is %s", fieldName, err.Tag()))
	}

	// Return the error messages in JSON format
	c.JSON(http.StatusBadRequest, gin.H{
		"error": errorMessages,
	})
}

// ValidateEmail validates the format and syntax of an email address using an external email verification package.
//
// Parameters:
// - email: The email address to be validated.
//
// Returns:
//   - An error if the email is invalid (either due to syntax or any other issues during verification).
//     Returns ErrInvalidEmailAddress if the syntax is not valid.
func ValidateEmail(email string) error {
	// Create a new verifier instance from the external package
	verifier := emailverifier.NewVerifier()

	// Verify the email address
	ret, err := verifier.Verify(email)
	if err != nil {
		// Return an error if the verification fails
		return fmt.Errorf("verify email address failed, error is: %v", err)
	}

	// If the syntax of the email is not valid, return an error
	if !ret.Syntax.Valid {
		return ErrInvalidEmailAddress
	}

	// Return nil if the email is valid
	return nil
}

// ValidatePhoneNumber validates a phone number's format and checks if it's valid using the phonenumbers package.
//
// Parameters:
// - phoneNum: The phone number to be validated (in any format).
//
// Returns:
// - The formatted phone number in E.164 format if the number is valid.
// - An error if the phone number is invalid or could not be parsed.
func ValidatePhoneNumber(phoneNum string) (string, error) {
	// Match a basic phone number format (starts with + or digits only)
	re := regexp.MustCompile(`^\+?[0-9]+$`)
	if !re.MatchString(phoneNum) {
		// Return error if the phone number format is incorrect
		return "", ErrInvaliPhoneNumber
	}

	// Parse the phone number using the phonenumbers library (assumes Uzbekistan as default country)
	num, err := phonenumbers.Parse(phoneNum, "UZ")
	if err != nil {
		// Return error if the phone number cannot be parsed
		return "", err
	}

	// Check if the phone number is valid
	if !phonenumbers.IsValidNumber(num) {
		// Return error if the phone number is not valid
		return "", ErrInvaliPhoneNumber
	}

	// Return the formatted phone number (E.164 format)
	return phonenumbers.Format(num, phonenumbers.E164), nil
}

// ValidatePassword validates the password based on specific criteria such as length, case sensitivity, numbers, and special characters.
//
// Parameters:
// - password: The password to be validated.
// - confirmPassword: A confirmation password to check if it matches the original password.
//
// Returns:
// - An error if the password doesn't meet the required criteria or if the confirmation password doesn't match the original.
func ValidatePassword(password, confirmPassword string) error {
	// Check if the password and confirm password match
	if password != confirmPassword {
		// Return an error if passwords don't match
		return ErrMismatchedPassword
	}

	// Validate password using various rules (length, case, numbers, special characters)
	err := validation.Validate(password,
		validation.Length(8, 0).Error("password must be at least 8 characters long"),                                                 // Minimum length of 8
		validation.Match(regexp.MustCompile(`[a-z]`)).Error("password must contain at least one lowercase letter"),                   // Lowercase letter required
		validation.Match(regexp.MustCompile(`[A-Z]`)).Error("password must contain at least one uppercase letter"),                   // Uppercase letter required
		validation.Match(regexp.MustCompile(`\d`)).Error("password must contain at least one number"),                                // At least one number required
		validation.Match(regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)).Error("password must contain at least one special character"), // Special character required
	)

	// Return any validation error if the password doesn't meet the criteria
	if err != nil {
		return err
	}

	// Return nil if the password is valid
	return nil
}
