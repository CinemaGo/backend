package handlers

import (
	"cinemaGo/backend/api/helpers"
	"cinemaGo/backend/internal/services"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	users services.UserServiceInterface
}

func NewUsersHandler(service services.UserServiceInterface) *UsersHandler {
	return &UsersHandler{users: service}
}

func (service *UsersHandler) SignUp(c *gin.Context) {
	var newUser newUserForm

	if err := c.ShouldBindJSON(&newUser); err != nil {
		helpers.RespondWithValidationErrors(c, err, newUser)
		return
	}

	if err := helpers.ValidateEmail(newUser.Email); err != nil {
		if errors.Is(err, helpers.ErrInvalidEmailAddress) {
			helpers.ClientError(c, http.StatusBadRequest, "invalid email address")
			return
		}
		helpers.ServerError(c, err)
		return
	}

	validPhoneNumber, err := helpers.ValidatePhoneNumber(newUser.PhoneNumber)
	if err != nil {
		if errors.Is(err, helpers.ErrInvaliPhoneNumber) {
			helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	if err := helpers.ValidatePassword(newUser.Password, newUser.ConfirmPassword); err != nil {
		if errors.Is(err, helpers.ErrMismatchedPassword) {
			helpers.ClientError(c, http.StatusBadRequest, "password and confirm password must match")
			return
		}
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	err = service.users.InsertNew(newUser.Name, newUser.Surname, newUser.Email, validPhoneNumber, newUser.Password)
	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmail) {
			helpers.ClientError(c, http.StatusConflict, fmt.Sprintf("%v", err))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully! Please login now",
	})
}

func (service *UsersHandler) Login(c *gin.Context) {
	var userLogin userLoginForm

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		helpers.RespondWithValidationErrors(c, err, userLogin)
		return
	}

	userID, userRole, err := service.users.UserAuthentication(userLogin.Email, userLogin.Password)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("%v", err))
		} else if errors.Is(err, services.ErrUserInvalidCredentials) {
			helpers.ClientError(c, http.StatusUnauthorized, fmt.Sprintf("%v", err))
		} else {
			helpers.ServerError(c, err)
		}
		return
	}

	tokenString, err := helpers.CreateAndSignJWT(userID, userRole)
	if err != nil {
		helpers.ServerError(c, err)
		return
	}

	helpers.SetCookie(c, tokenString)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Login complete! Explore what's new!",
	})
}

func (service *UsersHandler) UserProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	user_id := int(userID.(float64))

	userInfo, err := service.users.FetchUserInformations(user_id)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("%v", err))
		} else {
			helpers.ServerError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userInfo": userInfo,
	})
}

func (service *UsersHandler) UpdateUserProfile(c *gin.Context) {
	var userInfoUpdate userInfoUpdateFrom

	if err := c.ShouldBindJSON(&userInfoUpdate); err != nil {
		helpers.RespondWithValidationErrors(c, err, userInfoUpdate)
		return
	}

	validPhoneNumber, err := helpers.ValidatePhoneNumber(userInfoUpdate.PhoneNumber)
	if err != nil {
		if errors.Is(err, helpers.ErrInvaliPhoneNumber) {
			helpers.ClientError(c, http.StatusBadRequest, err.Error())
		} else {
			helpers.ServerError(c, err)
		}
		return
	}

	userID, _ := c.Get("userID")
	user_id := int(userID.(float64))

	err = service.users.UpdateUserInformations(user_id, userInfoUpdate.Name, userInfoUpdate.Surname, validPhoneNumber)
	if err != nil {
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
	})
}

func (service *UsersHandler) Logout(c *gin.Context) {
	c.SetCookie("u_auth", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "You are now logged out. Have a great day!",
	})
}
